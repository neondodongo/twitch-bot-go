package twitch

import "errors"

var (
	ErrInvalidConfig error = errors.New("twitch controller config is invalid")
	ErrPing          error = errors.New("failed to respond to Twitch ping")
	ErrNilConnection error = errors.New("the operator's network connection interface was nil")
	ErrNetworkRead   error = errors.New("failed to read from Twitch connection")
)
