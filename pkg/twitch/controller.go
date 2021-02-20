package twitch

import (
	"fmt"
	"net"
	"strings"
)

type Controller interface {
	/* Connection */

	Connect() error
	Disconnect()

	ObserveChat()
	JoinChannel() error
	Ping() error
	Write(msg string) error
	// SendMessage(msg string) error

	AddCommand(name string, cmd Command)
	GetCommand(name string) (Command, error)
	ListCommands() []string

	/* Getters */

	Channel() string
	Username() string
	ChatObserver() <-chan EventInfo
}

type operator struct {
	channel      string
	endpoint     string
	port         string
	creds        creds
	conn         net.Conn
	chatObserver chan EventInfo
	cmdRegistry  CommandRegistry
}

type creds struct {
	token    string
	username string
}

// New constructs a new Controller instance with the provided Controller Config, validating and sanitizing the required configuration values
func New(cfg Config) (Controller, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("failed construct new Twitch controller; %w", err)
	}

	return &operator{
		channel:  cfg.Channel,
		endpoint: cfg.URL,
		port:     cfg.Port,
		creds: creds{
			token:    cfg.OAuthToken,
			username: cfg.Username,
		},
		chatObserver: make(chan EventInfo),
		cmdRegistry:  make(map[string]Command),
	}, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Channel = strings.ToLower(strings.TrimSpace(cfg.Channel)); cfg.Channel == "" {
		return fmt.Errorf("%w; channel cannot be empty or whitespace only", ErrInvalidConfig)
	}

	if cfg.OAuthToken = strings.TrimSpace(cfg.OAuthToken); cfg.OAuthToken == "" {
		return fmt.Errorf("%w; oauth token cannot be empty or whitespace only", ErrInvalidConfig)
	}

	if cfg.Port = strings.TrimSpace(cfg.Port); cfg.Port == "" {
		return fmt.Errorf("%w; port cannot be empty or whitespace only", ErrInvalidConfig)
	}

	if cfg.URL = strings.TrimSpace(cfg.URL); cfg.URL == "" {
		return fmt.Errorf("%w; url cannot be empty or whitespace only", ErrInvalidConfig)
	}

	if cfg.Username = strings.ToLower(strings.TrimSpace(cfg.Username)); cfg.Username == "" {
		return fmt.Errorf("%w; username cannot be empty or whitespace only", ErrInvalidConfig)
	}

	return nil
}

func (op *operator) Channel() string {
	return op.channel
}

func (op *operator) Username() string {
	return op.creds.username
}

func (op *operator) ChatObserver() <-chan EventInfo {
	return op.chatObserver
}
