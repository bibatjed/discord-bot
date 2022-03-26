package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discord, err := discordgo.New("Bot " + os.Getenv("TWITCH_TOKEN"))
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

/**
ADD More commands
 */
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

	if cleanContent == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if cleanContent == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
