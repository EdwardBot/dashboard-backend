package database

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error

var mongoOnce sync.Once

const (
	DB            = "edward"
	USERS         = "users"
	SESSIONS      = "sessions"
	GUILDS        = "guilds"
	GUILD_CONFIGS = "guild-configs"
	WALLETS       = "wallets"
)

func Connect() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI") + "/")

	log.Println(clientOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to the database!\n%s", err.Error())
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

func GetGuilds() *mongo.Collection {
	return GetInstance().Database(DB).Collection(GUILDS)
}

func GetGuildConfigs() *mongo.Collection {
	return GetInstance().Database(DB).Collection(GUILD_CONFIGS)
}

func GetWalletCollection() *mongo.Collection {
	return GetInstance().Database(DB).Collection(WALLETS)
}
