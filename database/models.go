package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

var (
	userCollection    *mongo.Collection
	sessionCollection *mongo.Collection
	guildCollection   *mongo.Collection
	gConfCollection   *mongo.Collection
)

func Init() {
	userCollection = GetUsers()
	sessionCollection = GetSessions()
	guildCollection = GetGuilds()
	gConfCollection = GetGuildConfigs()
}

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        int64              `bson:"user_id"`
	JoinedAt      time.Time          `bson:"joined_at"`
	UserName      string             `bson:"username"`
	Discriminator string             `bson:"discriminator"`
	AvatarID      string             `bson:"avatar"`
	Guilds        []int64            `bson:"guilds"`
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       int64              `bson:"user_id"`
	AccessToken  string             `bson:"access_token"`
	RefreshToken string             `bson:"refresh_token"`
	RefreshedAt  time.Time          `bson:"refreshed_at"`
	ExpiresIn    int64              `bson:"expires_in"`
	SessionId    int32              `bson:"session_id"`
}

type Guild struct {
	ID         primitive.ObjectID `bson:"_id"`
	GuildID    int64              `bson:"guild_id"`
	GID 	   string             `bson:"gid"`
	Name       string             `bson:"name"`
	Icon       string             `bson:"icon"`
	HasBot     bool               `bson:"has_bot"`
	OwnerId    int64              `bson:"owner_id"`
	HasPremium bool               `bson:"has_premium"`
}

type GuildConfig struct {
	ID           primitive.ObjectID `bson:"_id"`
	AllowLogging bool               `bson:"allowLogging"`
	AllowWelcome bool               `bson:"allowWelcome"`
	BotAdmins    []string           `bson:"botAdmins"`
	GuildId      string             `bson:"guildId"`
	JoinedAt     int64              `bson:"joinedAt"`
}

func FindUser(id int64) (User, error) {
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

func FindSession(id int32) (Session, error) {
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

func FindGuild(id int64) (Guild, error) {
	var result Guild
	err := guildCollection.FindOne(context.TODO(), bson.D{
		primitive.E{
			Key:   "guild_id",
			Value: id,
		},
	}).Decode(&result)
	if err != nil {
		return Guild{}, err
	}
	return result, nil
}

func FindGuilds(userId int64) ([]Guild, error) {
	user, _ := FindUser(userId)

	guilds := make([]Guild, len(user.Guilds))

	for i, guild := range user.Guilds {
		g, _ := FindGuild(guild)
		g.GID  = strconv.FormatInt(g.GuildID, 10)
		guilds[i] = g
	}
	return guilds, nil
}

func FindGConf(id string) (GuildConfig, error) {
	var result GuildConfig
	err := gConfCollection.FindOne(context.TODO(), bson.D{
		primitive.E{
			Key:   "guildId",
			Value: id,
		},
	}).Decode(&result)
	if err != nil {
		return GuildConfig{}, err
	}
	return result, nil
}

func (u *User) Save() error {
	_, err := userCollection.InsertOne(context.TODO(), *u)
	return err
}

func (s *Session) Save() error {
	_, err := sessionCollection.InsertOne(context.TODO(), *s)
	return err
}

func (g *Guild) Save() error {
	_, err := guildCollection.InsertOne(context.TODO(), *g)
	return err
}

func (gc *GuildConfig) Save() error {
	_, err := gConfCollection.InsertOne(context.TODO(), *gc)
	return err
}

func (g *Guild) Update(id primitive.ObjectID) error {
	_, err := guildCollection.UpdateOne(context.TODO(), primitive.E{
		Key:   "_id",
		Value: id,
	}, *g)
	return err
}

func (gc *GuildConfig) Update(id primitive.ObjectID) error {
	_, err := gConfCollection.UpdateOne(context.TODO(), primitive.E{
		Key:   "_id",
		Value: id,
	}, *gc)
	return err
}

func (s *Session) Update(id primitive.ObjectID) error {
	_, err := sessionCollection.UpdateOne(context.TODO(), primitive.E{
		Key:   "_id",
		Value: id,
	}, *s)
	return err
}
