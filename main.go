package main

import (
	"twitch-notifier/config"
	"twitch-notifier/discord"
)

func init() {
	config.InitializeConfig()
}
func main() {
	discord.StartDiscord()
}
