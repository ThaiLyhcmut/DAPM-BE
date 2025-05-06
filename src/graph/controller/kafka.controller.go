package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	protoKafka "ThaiLy/proto/kafka"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

// UserChannelsMap manages user channels with mutex for thread safety
type UserChannelsMap struct {
	mu       sync.RWMutex
	channels map[int32][]chan *model.Device
}

// Global channels map with mutex protection
var userChannelsMap = UserChannelsMap{
	channels: make(map[int32][]chan *model.Device),
}

// Add adds a channel to the map
func (ucm *UserChannelsMap) Add(id int32, ch chan *model.Device) {
	ucm.mu.Lock()
	defer ucm.mu.Unlock()
	ucm.channels[id] = append(ucm.channels[id], ch)
}

// Remove removes a channel from the map
func (ucm *UserChannelsMap) Remove(id int32, ch chan *model.Device) {
	ucm.mu.Lock()
	defer ucm.mu.Unlock()

	channels := ucm.channels[id]
	for i, c := range channels {
		if c == ch {
			// Remove efficiently without preserving order
			channels[i] = channels[len(channels)-1]
			ucm.channels[id] = channels[:len(channels)-1]
			break
		}
	}

	// Clean up empty slices
	if len(ucm.channels[id]) == 0 {
		delete(ucm.channels, id)
	}
}

// Broadcast sends a message to all channels for a specific user
func (ucm *UserChannelsMap) Broadcast(id int32, device *model.Device) {
	ucm.mu.RLock()
	defer ucm.mu.RUnlock()

	for _, ch := range ucm.channels[id] {
		// Non-blocking send to prevent slow consumers from blocking the broadcast
		select {
		case ch <- device:
		default:
			// Channel is full or closed, will be cleaned up later
		}
	}
}

func (C *Controller) DeviceService(ctx context.Context, id int32, turnOn bool) (*string, error) {
	// Extract claims from context
	claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	accountID, err := helper.ParseASE(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("error parsing account ID: %w", err)
	}

	// Get equipment details
	equipment, err := C.equipment.CheckEquipment(id)
	if err != nil {
		return nil, fmt.Errorf("error getting equipment by ID: %w", err)
	}

	// Verify ownership
	checkEquipmentByAccountId, err := C.equipment.CheckHome(accountID, equipment.HomeId)
	if err != nil || checkEquipmentByAccountId == nil {
		return nil, fmt.Errorf("equipment does not belong to your home")
	}

	// Update equipment status in database
	if _, err = C.equipment.ChangeTurnOnEquipment(id, turnOn); err != nil {
		return nil, fmt.Errorf("failed to update equipment status: %w", err)
	}

	// Send to Kafka
	in := &protoKafka.DeviceRequest{
		Id:        id,
		TurnOn:    turnOn,
		AccountId: accountID,
	}

	res, err := C.kafka.DeviceService(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("kafka service error: %w", err)
	}

	return &res.Message, nil
}

// Single Kafka reader instance for all connections
var (
	kafkaReaderOnce sync.Once
	kafkaReader     *kafka.Reader
)

// getKafkaReader returns a singleton kafka reader
func getKafkaReader() *kafka.Reader {
	kafkaReaderOnce.Do(func() {
		kafkaReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKER")},
			Topic:   os.Getenv("DEVICE_TOGGLE_TOPIC"),
			GroupID: os.Getenv("GROUP_ID"),
			// Add performance settings
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
	})
	return kafkaReader
}

// startKafkaConsumer starts a single Kafka consumer for all connections
func (C *Controller) startKafkaConsumer(ctx context.Context) {
	reader := getKafkaReader()

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Kafka read error:", err)
				continue
			}

			// Process message
			parts := strings.Split(string(msg.Value), "|")
			if len(parts) != 3 {
				log.Println("Invalid message format:", string(msg.Value))
				continue
			}

			deviceID, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Println("Invalid device ID:", parts[0])
				continue
			}

			turnOn := parts[1] == "true"

			accountID, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				log.Println("Invalid account ID:", parts[2])
				continue
			}

			// Broadcast to all relevant channels
			device := &model.Device{ID: int32(deviceID), TurnOn: turnOn}
			userChannelsMap.Broadcast(int32(accountID), device)
		}
	}()
}

// Initialize Kafka consumer once
var initConsumerOnce sync.Once

func (C *Controller) DeviceStatusUpdated(ctx context.Context) (<-chan *model.Device, error) {
	// Start Kafka consumer if not already started
	initConsumerOnce.Do(func() {
		C.startKafkaConsumer(ctx)
	})

	claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("could not retrieve claims from context")
	}

	accountID, err := helper.ParseASE(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("error parsing ID: %w", err)
	}

	// Create buffered channel to prevent blocking
	ch := make(chan *model.Device, 100)

	// Add channel to map
	userChannelsMap.Add(accountID, ch)

	// Clean up channel when context is done
	go func() {
		<-ctx.Done()
		userChannelsMap.Remove(accountID, ch)
		close(ch)
	}()

	return ch, nil
}
