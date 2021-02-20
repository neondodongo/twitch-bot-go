package main

import (
	"twitch-bot-go/internal/client/okashi"
	"twitch-bot-go/pkg/twitch"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Booting up the bot")

	okashiBot, err := okashi.New(okashi.Config{
		Twitch: twitch.Config{
			OAuthToken: "auth-token-here-pls",
			Channel:    "neondodongo",
			Username:   "okashiibot",
			URL:        "irc.chat.twitch.tv",
			Port:       "6667",
		},
		TwitchVIPs: []okashi.TwitchVIP{
			{
				Username: "neondodongo",
				Messages: []string{
					"Hello father <3",
					"I missed you, dad!",
					"Daddy!!!",
				},
			},
			{
				Username: "okashiboxx",
				Messages: []string{
					"Hey there little lady ;)",
					"What's cookin' good lookin'?",
					"It's BAE <3",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	if err := okashiBot.StartTwitchBot(); err != nil {
		log.Fatal().Err(err).Msg("an error occurred starting OkashiBot")
	}
}
