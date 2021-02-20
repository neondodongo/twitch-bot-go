package okashi

// VIPRegistry is a map of specific usernames and their custom messages
type TwitchVIPRegistry map[string]TwitchVIP

type TwitchVIP struct {
	Username string   `json:"username"`
	Messages []string `json:"messages"`
	Greeted  bool     `json:"greeted"`
}

func initializeTwitchVIPRegistry(vips []TwitchVIP) map[string]TwitchVIP {
	var vipReg TwitchVIPRegistry = make(map[string]TwitchVIP)
	for _, vip := range vips {
		vipReg[vip.Username] = vip
	}

	return vipReg
}
