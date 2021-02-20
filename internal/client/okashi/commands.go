package okashi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"twitch-bot-go/pkg/twitch"
)

// Regex for parsing PRIVMSG strings.
//
// First matched group is the user's name and the second matched group is the content of the
// user's message.
var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)

// Regex for parsing user commands, from already parsed PRIVMSG strings.
//
// First matched group is the command name and the second matched group is the argument for the
// command.
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

func (bkr *broker) sayHello() error {
	defer bkr.clearCmdIssuer()

	username := bkr.cmdIssuer
	if username = strings.TrimSpace(username); username == "" {
		return fmt.Errorf("cannot say hello if I don't know the username")
	}

	if err := bkr.twitch.Write(fmt.Sprintf("@%s HeyGuys <3", username)); err != nil {
		return fmt.Errorf("failed to say hello; %w", err)
	}

	return nil
}

func (bkr *broker) greetVIP(vipUsername string) error {
	vip, ok := bkr.twitchVIPs[vipUsername]
	if !ok {
		return fmt.Errorf("user %s is not registered as a VIP", vipUsername)
	}

	// if VIP has already been greeted, don't greet them anymore
	if vip.Greeted {
		return nil
	}

	msg := ""

	if len(vip.Messages) > 1 {
		msg = vip.Messages[bkr.rand.Intn(len(vip.Messages))]
	} else if len(vip.Messages) == 0 {
		msg = vip.Messages[0]
	} else {
		return fmt.Errorf("no custom greeting detected for VIP %s", vipUsername)
	}

	if err := bkr.twitch.Write(fmt.Sprintf("@%s %s", vipUsername, msg)); err != nil {
		return fmt.Errorf("failed to send VIP greeting; %w", err)
	}

	// set greeted flag to true
	// TODO: this is a dirty trick, flag will only reset if the bot is reset
	vip.Greeted = true

	bkr.twitchVIPs[vipUsername] = vip

	return nil
}

func (bkr *broker) uwu() error {
	defer bkr.clearCmdIssuer()

	if err := bkr.twitch.Write("<(*~u w u~*)>"); err != nil {
		return fmt.Errorf("failed to send uwu message; %w", err)
	}

	return nil
}

func (bkr *broker) followage() error {
	defer bkr.clearCmdIssuer()

	username := bkr.cmdIssuer

	if username = strings.TrimSpace(username); username == "" {
		return fmt.Errorf("cannot get followage with an empty or whitespace-only username value")
	}

	if username == bkr.twitch.Channel() {
		return fmt.Errorf("a user cannot follow themself")
	}

	url := fmt.Sprintf("https://beta.decapi.me/twitch/followage/%s/%s", bkr.twitch.Channel(), username)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create http request for followage for user %s; %w", username, err)
	}

	res, err := bkr.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http request to obtain followage for user %s; %w", username, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected a 200 status code, but got %d", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to convert response body to bytes; %w", err)
	}

	if err := bkr.twitch.Write(fmt.Sprintf("%s has been following %s for %s", username, bkr.twitch.Channel(), string(bodyBytes))); err != nil {
		return fmt.Errorf("failed to write followage results to channel %s; %w", bkr.twitch.Channel(), err)
	}

	return nil
}

func (bkr *broker) subscribeMessage() twitch.Command {
	return func() error {
		if err := bkr.twitch.Write("Hello cuties! Don't forget to use your prime subscription and join the Snackpack okashi2Snacktime <3 xoxoxo"); err != nil {
			return fmt.Errorf("failed to write subscribe message to channel %s; %w", bkr.twitch.Channel(), err)
		}

		return nil
	}
}

func (bkr *broker) givawayMessage() error {
	if err := bkr.twitch.Write("Like what you're seeing? Give this channel a follow and turn on notifications. Once we hit 300 followers, Okashi will be hosting a $25 gaming gift card giveaway. ;)"); err != nil {
		return fmt.Errorf("failed to write giveaway message to channel %s; %w", bkr.twitch.Channel(), err)
	}

	return nil
}

func (bkr *broker) listCommands() error {
	username := bkr.cmdIssuer

	cmds := bkr.twitch.ListCommands()
	msg := ""
	for _, name := range cmds {
		msg += fmt.Sprintf("%s, ", name)
	}

	msg = strings.TrimSuffix(msg, ", ")

	if err := bkr.twitch.Write(fmt.Sprintf("@%s %s -- be sure to use the '!' prefix! Kappa", username, msg)); err != nil {
		return fmt.Errorf("failed to write giveaway message to channel %s; %w", bkr.twitch.Channel(), err)
	}

	return nil
}

func (bkr *broker) initCommandRegistry() {
	bkr.twitch.AddCommand("commands", bkr.listCommands)
	bkr.twitch.AddCommand("followage", bkr.followage)
	bkr.twitch.AddCommand("hello", bkr.sayHello)
	bkr.twitch.AddCommand("uwu", bkr.uwu)
}

func (bkr *broker) clearCmdIssuer() {
	bkr.cmdIssuer = ""
}
