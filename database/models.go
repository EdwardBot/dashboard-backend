package database

import (
	"gorm.io/gorm"
)

type User struct {
	UserID        int64  `json:"uid,omitempty" gorm:"index,primaryKey"`
	UID           string `json:"user_id,omitempty" gorm:"-"`
	UserName      string `json:"user_name,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	AvatarID      string `json:"avatar_id,omitempty"`
	Guilds        string `json:"guilds,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	JoinedAt      int64  `json:"joined_at,omitempty" gorm:"autoCreateTime"`
}

type Session struct {
	gorm.Model
	UserID       int64  `json:"user_id,omitempty" gorm:"index"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshedAt  int64  `json:"refreshed_at" gorm:"autoUpdateTime:milli"`
	ExpiresIn    uint64 `json:"expires_in,omitempty"`
}

type Guild struct {
	GuildID    uint64 `json:"-" gorm:"index;primaryKey;column:gid"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	HasBot     bool   `json:"has_bot"`
	OwnerId    int64  `json:"owner_id"`
	HasPremium bool   `json:"has_premium"`
	ID         string `json:"guild_id" gorm:"-"`
	Permission string `json:"permission" gorm:"-"`
}

type GuildConfig struct {
	BAs             string   `json:"-" gorm:"column:BotAdmins;type:integer[]"`
	GuildId         string   `json:"guild_id" gorm:"index;column:GuildId;primaryKey"`
	JoinedAt        string   `json:"joined_at" gorm:"column:JoinedAt"`
	Wch             uint64   `json:"-" gorm:"column:JoinChannel"`
	LCh             uint64   `json:"-" gorm:"column:LeaveChannel"`
	LoCh            uint64   `json:"-" gorm:"column:LogChannel"`
	WelcomeChannel  string   `json:"welcome_channel" gorm:"-"`
	LeaveChannel    string   `json:"leave_channel" gorm:"-"`
	LogChannel      string   `json:"log_channel" gorm:"-"`
	AllowedFeatures uint8    `json:"allowed_features" gorm:"column:AllowedFeatures"`
	BotAdmins       []string `json:"bot_admins" gorm:"-"`
}

type Wallet struct {
	ID       int32  `json:"id" gorm:"primaryKey;column:id"`
	GID      int64  `json:"-" gorm:"column:guild;index"`
	GuildId  string `json:"guild_id" gorm:"-"`
	UID      int64  `json:"-" gorm:"column:userid;index"`
	UserId   string `json:"user_id" gorm:"-"`
	Balance  int32  `json:"balance" gorm:"column:balance"`
	Xp       int32  `json:"xp" gorm:"column:xd"`
	Level    int32  `json:"level" gorm:"column:level"`
	Messages int32  `json:"messages" gorm:"column:messages"`
}

type Channel struct {
	ID    int64  `json:"id" gorm:"index"`
	Name  string `json:"name"`
	Guild int64  `json:"guild" gorm:"index"`
}

type Permissions struct {
	Guild       int64   `json:"-"`
	GuildId     string  `json:"guild_id" gorm:"-"`
	User        int64   `json:"-"`
	UserId      string  `json:"user_id" gorm:"-"`
	Perms       float64 `json:"-"`
	Permissions string  `json:"permissions" gorm:"-"`
}
