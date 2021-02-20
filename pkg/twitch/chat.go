package twitch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strings"

	"github.com/rs/zerolog/log"
)

// EventInfo represents an incoming message from the Twitch IRC connection(s)
type EventInfo struct {
	Message string
	Err     error
}

func (op *operator) ObserveChat() {
	tp := textproto.NewReader(bufio.NewReader(op.conn))

	go func() {
		for {
			line, err := tp.ReadLine()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Warn().Err(err).Msg("error occurred, attempting to reconnect")

					if err := op.Connect(); err != nil {
						op.chatObserver <- EventInfo{Err: fmt.Errorf("reconnect attempt failed; %w", err)}
						return
					}

					if err := op.JoinChannel(); err != nil {
						op.chatObserver <- EventInfo{Err: fmt.Errorf("failed to join channel after reconnect; %w", err)}
						return
					}
				}

				op.chatObserver <- EventInfo{Err: fmt.Errorf("%s; %w", ErrNetworkRead, err)}

				continue
			}

			if strings.Contains(line, TwitchPingMsg) {
				if err := op.Ping(); err != nil {
					op.chatObserver <- EventInfo{Err: err}
				}

				continue
			}

			op.chatObserver <- EventInfo{Message: line}
		}
	}()
}

func (op *operator) Write(msg string) error {
	if msg == "" {
		return fmt.Errorf("cannot write a message that is empty or whitespace-only")
	}

	if _, err := op.conn.Write([]byte(fmt.Sprintf("%s #%s :%s\r\n", CmdPrefixPrivateMsg, op.channel, msg))); err != nil {
		return fmt.Errorf("failed to write message to channel %s; %w", op.channel, err)
	}

	return nil
}
