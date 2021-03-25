package database

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error

var mongoOnce sync.Once

const (
	DB       = "edward"
	USERS    = "users"
	SESSIONS = "sessions"
)

func Connect() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI") + "/" + DB)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		clientInstanceError = err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		clientInstanceError = err
	}
	clientInstance = client
	return clientInstanceError
}

func GetInstance() *mongo.Client {
	return clientInstance
}

func GetUsers() *mongo.Collection {
	return GetInstance().Database(DB).Collection(USERS)
}

func GetSessions() *mongo.Collection {
	return GetInstance().Database(DB).Collection(SESSIONS)
}
