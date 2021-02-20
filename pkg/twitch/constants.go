package twitch

const (
	// Command Prefixes

	CmdPrefixPassword    string = "PASS"
	CmdPrefixNickname    string = "NICK"
	CmdPrefixJoinChannel string = "JOIN"
	CmdPrefixPrivateMsg  string = "PRIVMSG"

	// Known incoming messages

	TwitchPingMsg string = "PING :tmi.twitch.tv"
)
