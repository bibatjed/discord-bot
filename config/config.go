package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TwitchClientID     string
	TwitchClientSecret string
	DiscordToken       string
}

func InitializeConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TwitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	TwitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	DiscordToken := os.Getenv("DISCORD_TOKEN")

	return &Config{TwitchClientID: TwitchClientID, TwitchClientSecret: TwitchClientSecret, DiscordToken: DiscordToken}
}
