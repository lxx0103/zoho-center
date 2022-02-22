package auth

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ApiDomain   string `json:"api_domain"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
