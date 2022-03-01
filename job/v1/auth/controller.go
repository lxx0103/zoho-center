package auth

import "time"

func GetCode() (string, error) {
	duration := time.Duration(3) * time.Second
	time.Sleep(duration)
	authService := NewAuthService()
	token, err := authService.GetToken()
	if err != nil {
		return "", err
	}
	return token, nil

}
