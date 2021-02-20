package okashi

import "twitch-bot-go/pkg/twitch"

type Config struct {
	Twitch     twitch.Config `json:"twitch"`
	TwitchVIPs []TwitchVIP   `json:"twitch_vips"`
}
