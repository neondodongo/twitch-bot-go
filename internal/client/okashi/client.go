package okashi

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"twitch-bot-go/pkg/twitch"
)

type Client interface {
	StartTwitchBot() error
	Twitch() twitch.Controller
}

type broker struct {
	twitch     twitch.Controller
	http       *http.Client
	rand       *rand.Rand
	twitchVIPs TwitchVIPRegistry
	cmdIssuer  string
}

func New(cfg Config) (Client, error) {
	twitchCtrl, err := twitch.New(cfg.Twitch)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize new Twitch bot with provided config; %w", err)
	}

	bkr := &broker{
		twitch:     twitchCtrl,
		http:       http.DefaultClient, // TODO: configure custom http client
		twitchVIPs: initializeTwitchVIPRegistry(cfg.TwitchVIPs),
		rand:       rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}

	bkr.initCommandRegistry()

	return bkr, nil
}

func (bkr *broker) StartTwitchBot() error {
	if err := bkr.twitch.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Twitch host; %w", err)
	}

	if err := bkr.twitch.JoinChannel(); err != nil {
		return fmt.Errorf("failed to join channel; %w", err)
	}

	bkr.observeChat()

	return nil
}

func (bkr *broker) Twitch() twitch.Controller {
	return bkr.twitch
}
