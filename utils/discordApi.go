package utils

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var (
	discordInstance *discordgo.Session
)

func InitDiscord() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error connecting to discord!\n%s", err.Error())
	}
	discordInstance = discord
}

func GetDiscordInstance() *discordgo.Session {
	return discordInstance
}

func Find(array []*discordgo.Role, check func(e *discordgo.Role) bool) *discordgo.Role {
	for i := range array {
		if check(array[i]) {
			return array[i]
		}
	}
	return nil
}
