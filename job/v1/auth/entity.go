package auth

import "time"

type Token struct {
	ID          int64     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	AccessToken string    `db:"access_token" json:"access_token"`
	ApiDomain   string    `db:"api_domain" json:"api_domain"`
	TokenType   string    `db:"token_type" json:"token_type"`
	ExpiresTime time.Time `db:"expires_time" json:"expires_time"`
}
