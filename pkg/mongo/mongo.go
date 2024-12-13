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

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient() (*Client, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is not set")
	}

	db := client.Database(dbName)
	log.Printf("Connected to MongoDB: %s, database: %s\n", mongoURI, dbName)

	return &Client{client: client, db: db}, nil
}

type Record struct {
	Key        string    `bson:"key" json:"key"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}

type FilteredRecords struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []Record `json:"records"`
}

func (m *Client) FetchData(startDate, endDate string, minCount, maxCount int) ([]Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("records")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Printf("Error parsing start date: %v\n", err)
		return nil, fmt.Errorf("invalid start date format")
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Printf("Error parsing end date: %v\n", err)
		return nil, fmt.Errorf("invalid end date format")
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": startTime,
					"$lte": endTime,
				},
			},
		},
		{
			"$addFields": bson.M{
				"totalCount": bson.M{"$sum": "$counts"},
			},
		},
		{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gte": minCount,
					"$lte": maxCount,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"key":        1,
				"createdAt":  1,
				"totalCount": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error fetching data from MongoDB: %v\n", err)
		return nil, fmt.Errorf("error fetching data from MongoDB: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}(cursor, ctx)

	var results []Record
	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			log.Printf("Error decoding record: %v\n", err)
			return nil, fmt.Errorf("error decoding record from mongo: %w", err)
		}
		results = append(results, record)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration: %v\n", err)
		return nil, fmt.Errorf("error during cursor iteration: %w", err)
	}

	return results, nil
}
