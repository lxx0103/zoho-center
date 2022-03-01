package auth

func GetCode() (string, error) {
	authService := NewAuthService()
	token, err := authService.GetToken()
	if err != nil {
		return "", err
	}
	return token, nil

}
