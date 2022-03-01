package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"zoho-center/core/config"
	"zoho-center/core/database"
)

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

type AuthService interface {
	GetToken() (string, error)
	// UpdateUser(int64, UserUpdate, int64) (*User, error)
	// //Role Management
	// GetRoleByID(int64) (*Role, error)
	// NewRole(RoleNew) (*Role, error)
	// GetRoleList(RoleFilter) (int, *[]Role, error)
	// UpdateRole(int64, RoleNew) (*Role, error)
}

func (s authService) GetToken() (string, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return "", errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)

	clientCode := config.ReadConfig("zoho.client_code")
	tokenNow, err := repo.GetTokenByCode(clientCode)
	if err != nil {
		return "", err
	}
	if time.Now().Local().Add(5 * time.Microsecond).Before(tokenNow.ExpiresTime) {
		return tokenNow.AccessToken, nil
	}
	url := config.ReadConfig("zoho.token_uri")
	refreshToken := config.ReadConfig("zoho.refresh_token")
	clientID := config.ReadConfig("zoho.client_id")
	clientSecret := config.ReadConfig("zoho.client_secret")
	redirectUri := config.ReadConfig("zoho.redirect_uri")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("refresh_token", refreshToken)
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("redirect_uri", redirectUri)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tokenStruct TokenResponse
	err = json.Unmarshal(body, &tokenStruct)
	if err != nil {
		return "", err
	}
	tokenNow.AccessToken = tokenStruct.AccessToken
	tokenNow.ExpiresTime = time.Now().Local().Add(time.Second * time.Duration(tokenStruct.ExpiresIn))
	tokenNow.ApiDomain = tokenStruct.ApiDomain
	tokenNow.TokenType = tokenStruct.TokenType
	err = repo.UpdateToken(tokenNow.ID, *tokenNow)
	if err != nil {
		return "", err
	}
	tx.Commit()
	return tokenNow.AccessToken, nil
}
