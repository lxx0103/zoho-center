package auth

func GetCode() (string, error) {
	tokenService := NewAuthService()
	token, err := tokenService.GetToken()
	if err != nil {
		return "", err
	}
	return token, nil

}
