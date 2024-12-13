package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoClient() (*DBClient, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	dbName := "development"
	db := client.Database(dbName)
	log.Printf("Connected to MongoDB: %s\n", mongoURI)

	return &DBClient{client: client, db: db}, nil
}

func (m *DBClient) FetchData(collectionName string) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching data from MongoDB, %v \n", err)
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error fetching data from MongoDB, %v \n", err)
		}
	}(cursor, ctx)

	var results []bson.M

	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("Error decoding data from MongoDB, %v \n", err)
		return nil, err
	}
	return results, nil
}
