package okashi

import (
	"strings"
	"twitch-bot-go/pkg/twitch"

	"github.com/rs/zerolog/log"
)

func (bkr *broker) observeChat() {
	// messageChan := make(chan string)

	go func() {
		// subscribeMessageCounter := 0
		// giveawayMessageCounter := 0
		// subscribeTimer := time.NewTimer(15 * time.Second)
		// giveawayTimer := time.NewTimer(20 * time.Second)

		// for {
		// 	select {
		// 	case <-messageChan:
		// 		subscribeMessageCounter++
		// 		giveawayMessageCounter++
		// 	case <-subscribeTimer.C:
		// 		if subscribeMessageCounter > 0 {
		// 			if err := bkr.SubscribeMessage(); err != nil {
		// 				log.Info().Err(err).Msg("failed to send subscribe message")
		// 			}
		// 		}

		// 		subscribeTimer.Reset(15 * time.Second)
		// 		subscribeMessageCounter = 0
		// 	case <-giveawayTimer.C:
		// 		if giveawayMessageCounter > 0 {
		// 			if err := bkr.GivawayMessage(); err != nil {
		// 				log.Info().Err(err).Msg("failed to send giveaway message")
		// 			}
		// 		}

		// 		giveawayTimer.Reset(20 * time.Second)
		// 		giveawayMessageCounter = 0
		// 	}
		// }
	}()

	bkr.twitch.ObserveChat()

	for event := range bkr.twitch.ChatObserver() {
		if event.Err != nil {
			log.Error().Err(event.Err).Msgf("error occurred while observing channel %s", bkr.twitch.Channel())
			continue
		}

		log.Info().Msgf("Event received %s", event.Message)

		// messageChan <- data.Message

		if messageMatch := msgRegex.FindStringSubmatch(event.Message); messageMatch != nil {
			chatUsername := messageMatch[1]
			msgType := messageMatch[2]

			if chatUsername == bkr.twitch.Username() {
				log.Debug().Msgf("this bot will not respond to it's own messages")
				continue
			}

			if err := bkr.greetVIP(chatUsername); err != nil {
				log.Info().Err(err).Msg("error occurred while greeting VIP")
			}

			if msgType == twitch.CmdPrefixPrivateMsg {
				msg := messageMatch[3]

				if cmdMatches := cmdRegex.FindStringSubmatch(msg); cmdMatches != nil {
					cmdName := strings.TrimSpace(cmdMatches[1])

					cmd, err := bkr.twitch.GetCommand(cmdName)
					if err != nil {
						log.Info().Err(err).Msg("failed to retrieve Twitch command")
						continue
					}

					if cmd == nil {
						continue
					}

					bkr.cmdIssuer = chatUsername

					if err := cmd(); err != nil {
						log.Error().Err(err).Msg("failed to execute Twitch command")
					}
				}
			}
		}
	}
}
