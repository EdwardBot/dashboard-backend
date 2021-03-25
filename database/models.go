package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	userCollection    *mongo.Collection
	sessionCollection *mongo.Collection
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        int64              `bson:"user_id"`
	JoinedAt      time.Time          `bson:"joined_at"`
	UserName      string             `bson:"username"`
	Discriminator int                `bson:"discriminator"`
	AvatarID      string             `bson:"avatar"`
	Guilds        []int64            `bson:"guilds"`
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       int64              `bson:"user_id"`
	AccessToken  string             `bson:"access_token"`
	RefreshToken string             `bson:"refresh_token"`
	RefreshedAt  time.Time          `bson:"refreshed_at"`
	ExpiresIn    int                `bson:"expires_in"`
	SessionId    int64              `bson:"session_id"`
}

func FindUser(id string) (User, error) {
	if userCollection == nil {
		userCollection = GetUsers()
	}
	var result User
	err := userCollection.FindOne(context.TODO(), bson.D{
		primitive.E{
			Key:   "user_id",
			Value: id,
		},
	}).Decode(&result)
	if err != nil {
		return User{}, err
	}
	return result, nil
}

func FindSession(id int64) (Session, error) {
	if sessionCollection == nil {
		sessionCollection = GetSessions()
	}
	var result Session
	err := sessionCollection.FindOne(context.TODO(), bson.D{
		primitive.E{
			Key:   "session_id",
			Value: id,
		},
	}).Decode(&result)
	if err != nil {
		return Session{}, err
	}
	return result, nil
}

func (u *User) Save() error {
	if userCollection == nil {
		userCollection = GetUsers()
	}
	_, err := userCollection.InsertOne(context.TODO(), u)
	return err
}

func (s *Session) Save() error {
	if sessionCollection == nil {
		sessionCollection = GetSessions()
	}
	_, err := sessionCollection.InsertOne(context.TODO(), s)
	return err
}
