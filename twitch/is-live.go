package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"twitch-notifier/config"
)

var client *helix.Client

type ChannelStatus struct {
	Status      string
	DisplayName string
	LiveStatus  string
}

func init() {
	result := config.InitializeConfig()
	var err error
	client, err = helix.NewClient(&helix.Options{
		ClientID:     result.TwitchClientID,
		ClientSecret: result.TwitchClientSecret,
	})
	if err != nil {
		panic(err)
	}

}

func IsLive(channel string) *ChannelStatus {
	resps, errs := client.RequestAppAccessToken([]string{"user:read:email"})
	if errs != nil {
		// handle error
	}

	// Set the access token on the client
	client.SetAppAccessToken(resps.Data.AccessToken)

	channelInformation, _ := client.SearchChannels(&helix.SearchChannelsParams{
		Channel: channel,
		First:   1,
	})

	if len(channelInformation.Data.Channels) <= 0 {
		return &ChannelStatus{Status: "not_found"}
	}

	isLive := &channelInformation.Data.Channels[0].IsLive

	var channelStatus = "Not Live"

	if *isLive {
		channelStatus = "Live"
	}

	return &ChannelStatus{LiveStatus: channelStatus, DisplayName: *(&channelInformation.Data.Channels[0].DisplayName), Status: "found"}
}
