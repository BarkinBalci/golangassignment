package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DBClient struct {
	client *mongo.Client
	db     *mongo.Database
}
