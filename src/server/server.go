package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"ThaiLy/graph/controller"
	"ThaiLy/graph/generated"
	"ThaiLy/graph/helper"
	"ThaiLy/graph/resolver"
	protoKafka "ThaiLy/proto/kafka"

	"ThaiLy/server/client"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/vektah/gqlparser/v2/ast"
)

func handleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	kafkaTopic := "device_status"
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   kafkaTopic,
	})
	defer kafkaWriter.Close()

	// Parse nội dung MQTT message
	var req protoKafka.DeviceRequest
	err := json.Unmarshal(msg.Payload(), &req)
	if err != nil {
		log.Printf("❌ Lỗi parse MQTT message: %v", err)
		return
	}

	// Format lại dữ liệu giống ToggleDevice
	message := fmt.Sprintf("%d|%t|%s", req.Id, req.TurnOn, req.AccountId)

	// Ghi vào Kafka
	err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", req.Id)),
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("❌ Lỗi ghi vào Kafka: %v", err)
	} else {
		log.Printf("✅ Ghi vào Kafka thành công: %s", message)
	}
}

func connectMQTT(broker, clientID, username, password, topic string) {
	fmt.Println(broker, clientID, username, password, topic)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Thêm username và password
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(handleMQTTMessage)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Lỗi kết nối MQTT: %v", token.Error())
	}
	log.Println("✅ Kết nối MQTT thành công!")

	// Đăng ký lắng nghe topic MQTT
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Lỗi đăng ký topic MQTT: %v", token.Error())
	}
	log.Printf("📩 Đang lắng nghe MQTT topic: %s", topic)
}

const defaultPort = "8081"

func main() {

	gin.SetMode(gin.ReleaseMode)

	// Tải biến môi trường
	godotenv.Load()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Chỉ cho phép origin cụ thể
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	go connectMQTT(os.Getenv("MQTT_BROKER"), "myClient", os.Getenv("MQTT_USER"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC"))
	auth, err := client.NewGRPCAuthClient(os.Getenv("SERVICE_AUTH"))
	if err != nil {
		log.Fatalf("client auth error %v", err)
	}
	equipment, err := client.NewGRPCEquipmentClient(os.Getenv("SERVICE_EQUIPMENT"))
	if err != nil {
		log.Fatalf("client equipment error %v", err)
	}
	kafka, err := client.NewGRPCKafkaClient(os.Getenv("SERVICE_KAFKA"))
	if err != nil {
		log.Fatalf("client equipment error %v", err)
	}
	ctrl := controller.NewController(auth, equipment, kafka)
	// Create GraphQL handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{Ctrl: ctrl}}))

	// Add GraphQL transports (Options, GET, POST)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})
	// Set query cache
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Gin routes foyground and query handler
	r.GET("/", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
	r.Any("/query", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
				Claims, err := helper.ParseJWT(token)
				if Claims == nil || err != nil {
					c.JSON(401, gin.H{"message": "Invalid Authorization header"})
					c.Abort()
					return
				}
				ctx := context.WithValue(c.Request.Context(), helper.Auth, Claims)
				c.Request = c.Request.WithContext(ctx)
			}
		}
	}, func(c *gin.Context) {
		strings.Contains(c.GetHeader("Upgrade"), "websocket")
		gin.WrapH(srv)(c)
	})

	// Run the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
