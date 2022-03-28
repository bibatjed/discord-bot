package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix/v2"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type ChannelStatus struct {
	DisplayName string
	IsLive      bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
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

	splitString := strings.Split(cleanContent, " ")

	if len(splitString) <= 0 {
		return
	}

	if splitString[0] == "is-live" {

		if len(splitString) < 2 {
			return
		}
		client, err := helix.NewClient(&helix.Options{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		})

		if err != nil {
			panic(err)
		}

		resps, errs := client.RequestAppAccessToken([]string{"user:read:email"})
		if errs != nil {
			// handle error
		}

		// Set the access token on the client
		client.SetAppAccessToken(resps.Data.AccessToken)

		channelInformation, _ := client.SearchChannels(&helix.SearchChannelsParams{
			Channel: splitString[1],
			First:   1,
		})

		if len(channelInformation.Data.Channels) <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Not found")
			return
		}

		isLive := &channelInformation.Data.Channels[0].IsLive

		var channelStatus = "Not Live"

		if *isLive {
			channelStatus = "Live"
		}

		s.ChannelMessageSend(m.ChannelID, channelInformation.Data.Channels[0].DisplayName+": "+channelStatus)
	}
}
