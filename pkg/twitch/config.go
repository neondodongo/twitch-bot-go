package twitch

type Config struct {
	OAuthToken string `json:"oauth_token"`
	Channel    string `json:"channel"` // must be lowercase
	Password   string `json:"password"`
	Port       string `json:"port"`
	URL        string `json:"url"`
	Username   string `json:"username"` // must be lowercase
}
