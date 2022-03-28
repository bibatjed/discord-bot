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

	resps, errs := client.RequestAppAccessToken([]string{"user:read:email"})
	if errs != nil {
		// handle error
	}

	//Set the access token on the client
	client.SetAppAccessToken(resps.Data.AccessToken)
}

func IsLive(channel string, retry int) *ChannelStatus {
	channelInformation, _ := client.SearchChannels(&helix.SearchChannelsParams{
		Channel: channel,
		First:   1,
	})

	if retry > 0 && channelInformation.StatusCode == 401 {
		resps, errs := client.RequestAppAccessToken([]string{"user:read:email"})
		if errs != nil {
			// handle error
		}
		client.SetAppAccessToken(resps.Data.AccessToken)
		retry -= 1
		return IsLive(channel, retry)
	}

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
