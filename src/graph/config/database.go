// config/mongo.go
package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB(uri string) {
	var err error
	MongoClient, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error creating MongoDB client: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB")
}

const collectionName = "connections"

// Lưu trạng thái kết nối của user
func SaveUserConnection(userID int32) error {
	collection := MongoClient.Database("websocket").Collection(collectionName)

	// Dữ liệu lưu vào MongoDB
	connection := bson.M{
		"userID":      userID,
		"status":      "connected",
		"connectedAt": time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), connection)
	if err != nil {
		log.Println("❌ Lỗi lưu kết nối:", err)
		return err
	}
	log.Println("✅ Đã lưu kết nối cho User:", userID)
	return nil
}

// Xóa trạng thái khi user ngắt kết nối
func RemoveUserConnection(userID int32) error {
	collection := MongoClient.Database("websocket").Collection(collectionName)

	_, err := collection.DeleteOne(context.Background(), bson.M{"userID": userID})
	if err != nil {
		log.Println("❌ Lỗi xóa kết nối:", err)
		return err
	}
	log.Println("❌ User ngắt kết nối:", userID)
	return nil
}

// IsUserConnected kiểm tra xem user có đang kết nối không
func IsUserConnected(userID int32) (bool, error) {
	collection := MongoClient.Database("websocket").Collection(collectionName)

	// Tìm userId trong MongoDB
	filter := bson.M{"userID": userID}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
