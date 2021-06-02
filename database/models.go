package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

var (
	userCollection           *mongo.Collection
	sessionCollection        *mongo.Collection
	guildCollection          *mongo.Collection
	gConfCollection          *mongo.Collection
	walletCollection         *mongo.Collection
	customCommandsCollection *mongo.Collection
)

func Init() {
	userCollection = GetUsers()
	sessionCollection = GetSessions()
	guildCollection = GetGuilds()
	gConfCollection = GetGuildConfigs()
	walletCollection = GetWalletCollection()
	customCommandsCollection = GetCustomCommandsCollection()
}

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        int64              `bson:"user_id"`
	JoinedAt      time.Time          `bson:"joined_at"`
	UserName      string             `bson:"username"`
	Discriminator string             `bson:"discriminator"`
	AvatarID      string             `bson:"avatar"`
	Guilds        []int64            `bson:"guilds"`
	PremiumType   int                `bson:"premium_type"`
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
	ID          primitive.ObjectID `bson:"_id"`
	GuildID     int64              `bson:"guild_id"`
	GID         string             `bson:"gid"`
	Name        string             `bson:"name"`
	Icon        string             `bson:"icon"`
	HasBot      bool               `bson:"has_bot"`
	OwnerId     int64              `bson:"owner_id"`
	HasPremium  bool               `bson:"has_premium"`
	Permissions string
}

type GuildConfig struct {
	ID           primitive.ObjectID `bson:"_id"`
	AllowLogging bool               `bson:"allowLogging"`
	AllowWelcome bool               `bson:"allowWelcome"`
	BotAdmins    []string           `bson:"botAdmins"`
	GuildId      string             `bson:"guildId"`
	JoinedAt     int64              `bson:"joinedAt"`
}

type Wallet struct {
	ID       primitive.ObjectID `bson:"_id"`
	GuildId  string             `bson:"guildId"`
	UserId   string             `bson:"userId"`
	GID      string
	UID      string
	Balance  int32 `bson:"balance"`
	Xp       int32 `bson:"xp"`
	Lvl      int32 `bson:"lvl"`
	Messages int32 `bson:"messages"`
}

type CustomCommand struct {
	ID       primitive.ObjectID `bson:"_id"`
	GuildId  string             `bson:"guildId"`
	Name     string             `bson:"name"`
	Response string             `bson:"response"`
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

func FindWallet(userId, guildId string) *Wallet {
	var wallet Wallet
	e := walletCollection.FindOne(context.TODO(), bson.D{
		primitive.E{
			Key:   "guildId",
			Value: guildId,
		},
		primitive.E{
			Key:   "userId",
			Value: userId,
		},
	}).Decode(&wallet)
	if e != nil {
		log.Println(e.Error())
	}
	return &wallet
}

var empty = make([]CustomCommand, 0, 0)

func FindCommands(guildId string) *[]CustomCommand {
	var cmd = make([]CustomCommand, 0, 0)
	c, e := customCommandsCollection.Find(context.TODO(), bson.D{
		primitive.E{
			Key:   "guildId",
			Value: guildId,
		},
	})
	if e != nil {
		return &empty
	}

	for c.Next(context.TODO()) {
		var tmp CustomCommand
		c.Decode(&tmp)
		cmd = append(cmd, tmp)
	}

	return &cmd
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
		g.GID = strconv.FormatInt(g.GuildID, 10)
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

func (cc *CustomCommand) Save() error {
	err := customCommandsCollection.FindOne(context.TODO(), bson.D{
		bson.E{
			Key:   "_id",
			Value: cc.ID,
		},
	})
	if err != nil {
		_, err := customCommandsCollection.InsertOne(context.TODO(), *cc)
		return err
	} else {
		_, err := customCommandsCollection.UpdateOne(context.TODO(), bson.D{
			bson.E{
				Key:   "_id",
				Value: cc.ID,
			},
		}, *cc)

		return err
	}
}

func FindCommand(guildId, name string) (*CustomCommand, error) {
	var res CustomCommand
	err := customCommandsCollection.FindOne(context.TODO(), bson.D{
		bson.E{
			Key:   "guildId",
			Value: guildId,
		},
		bson.E{
			Key:   "name",
			Value: name,
		},
	}).Decode(&res)
	return &res, err
}

func (cc *CustomCommand) Delete() error {
	_, err := customCommandsCollection.DeleteOne(context.TODO(), bson.D{
		bson.E{
			Key:   "_id",
			Value: cc.ID,
		},
	})
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

func (u *User) Update(id primitive.ObjectID) error {
	_, err := userCollection.UpdateOne(context.TODO(), primitive.E{
		Key:   "_id",
		Value: id,
	}, *u)
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
