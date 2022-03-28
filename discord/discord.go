package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"twitch-notifier/config"
	"twitch-notifier/twitch"
)

func StartDiscord() {
	result := config.InitializeConfig()
	discord, err := discordgo.New("Bot " + result.DiscordToken)
	defer discord.Close()

	if err != nil {
		fmt.Println("Can't initialize discord connection")
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var botID = "<@!" + s.State.User.ID + ">"
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Ignore message if bot is not mentioned
	if ok := strings.Contains(m.Content, botID); !ok {
		return
	}

	var cleanContent = strings.ReplaceAll(m.Content, botID, "")
	cleanContent = strings.TrimSpace(cleanContent)

	splitString := strings.Split(cleanContent, " ")

	if len(splitString) <= 0 {
		return
	}

	if splitString[0] == "is-live" {

		if len(splitString) < 2 {
			return
		}

		result := twitch.IsLive(splitString[1], 3)

		if result.Status == "not_found" {
			s.ChannelMessageSend(m.ChannelID, "Not found")
			return
		}

		s.ChannelMessageSend(m.ChannelID, result.DisplayName+": "+result.LiveStatus)
	}
}
