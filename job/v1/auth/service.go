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

// func (s *authService) VerifyWechatSignin(code string) (*WechatCredential, error) {
// 	var credential WechatCredential
// 	httpClient := &http.Client{}
// 	signin_uri := config.ReadConfig("Wechat.signin_uri")
// 	appID := config.ReadConfig("Wechat.app_id")
// 	appSecret := config.ReadConfig("Wechat.app_secret")
// 	uri := signin_uri + "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
// 	req, err := http.NewRequest("GET", uri, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, err := httpClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(body, &credential)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &credential, nil
// }

// func (s *authService) GetUserInfo(openID string) (*User, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	user, err := query.GetUserByOpenID(openID)
// 	if err != nil {
// 		if err.Error() != "sql: no rows in result set" {
// 			return nil, err
// 		}
// 		tx, err := db.Begin()
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer tx.Rollback()
// 		var newUser User
// 		newUser.Type = 2
// 		newUser.Identifier = openID
// 		repo := NewAuthRepository(tx)
// 		userID, err := repo.CreateUser(newUser)
// 		if err != nil {
// 			return nil, err
// 		}
// 		user, err = repo.GetUserByID(userID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tx.Commit()
// 	}
// 	return user, nil
// }

// func (s *authService) VerifyCredential(signinInfo SigninRequest) (*User, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	userInfo, err := query.GetUserByUserName(signinInfo.Identifier)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !checkPasswordHash(signinInfo.Credential, userInfo.Credential) {
// 		errMessage := "密码错误"
// 		return nil, errors.New(errMessage)
// 	}
// 	return userInfo, err
// }

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// func checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// func (s *authService) UpdateUser(userID int64, info UserUpdate, byUserID int64) (*User, error) {
// 	db := database.InitMySQL()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		msg := "事务开启错误" + err.Error()
// 		return nil, errors.New(msg)
// 	}
// 	defer tx.Rollback()
// 	repo := NewAuthRepository(tx)

// 	oldUser, err := repo.GetUserByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	byUser, err := repo.GetUserByID(byUserID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	byRole, err := repo.GetRoleByID(byUser.RoleID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if oldUser.RoleID != 0 {
// 		targetRole, err := repo.GetRoleByID(oldUser.RoleID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if byRole.Priority <= targetRole.Priority && userID != byUserID { //只能修改角色比自己优先级低的用户,或者用户自身
// 			msg := "你无法修改角色为" + targetRole.Name + "的用户"
// 			return nil, errors.New(msg)
// 		}
// 	}
// 	if info.RoleID != 0 {
// 		toRole, err := repo.GetRoleByID(info.RoleID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if byRole.Priority < toRole.Priority { //只能将目标修改为和自己同级的角色
// 			msg := "你无法将目标角色改为:" + toRole.Name
// 			return nil, errors.New(msg)
// 		}
// 		oldUser.RoleID = info.RoleID
// 	}
// 	if info.PositionID != 0 {
// 		oldUser.PositionID = info.PositionID
// 	}
// 	if info.Name != "" {
// 		oldUser.Name = info.Name
// 	}
// 	if info.Email != "" {
// 		oldUser.Email = info.Email
// 	}
// 	if info.Gender != "" {
// 		oldUser.Gender = info.Gender
// 	}
// 	if info.Birthday != "" {
// 		oldUser.Birthday = info.Birthday
// 	}
// 	if info.Phone != "" {
// 		oldUser.Phone = info.Phone
// 	}
// 	if info.Address != "" {
// 		oldUser.Address = info.Address
// 	}
// 	if info.Status != 0 {
// 		oldUser.Status = info.Status
// 	}
// 	err = repo.UpdateUser(userID, *oldUser, (*byUser).Identifier)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := repo.GetUserByID(userID)
// 	tx.Commit()
// 	return user, err
// }

// func (s *authService) GetRoleByID(id int64) (*Role, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	role, err := query.GetRoleByID(id)
// 	return role, err
// }

// func (s *authService) NewRole(info RoleNew) (*Role, error) {
// 	db := database.InitMySQL()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()
// 	repo := NewAuthRepository(tx)
// 	roleID, err := repo.CreateRole(info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	role, err := repo.GetRoleByID(roleID)
// 	tx.Commit()
// 	return role, err
// }

// func (s *authService) GetRoleList(filter RoleFilter) (int, *[]Role, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	count, err := query.GetRoleCount(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	list, err := query.GetRoleList(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	return count, list, err
// }

// func (s *authService) UpdateRole(roleID int64, info RoleNew) (*Role, error) {
// 	db := database.InitMySQL()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()
// 	repo := NewAuthRepository(tx)
// 	_, err = repo.UpdateRole(roleID, info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	role, err := repo.GetRoleByID(roleID)
// 	tx.Commit()
// 	return role, err
// }
