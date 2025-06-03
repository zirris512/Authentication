package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zirris512/Authentication/database"
)

type VerifiedUser struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type PostBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSecretKey() ([]byte, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return []byte{}, fmt.Errorf("environment variable SECRET_KEY not set")
	}
	return secretKey, nil
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	secretKey, err := GetSecretKey()
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secretKey, err := GetSecretKey()
		if err != nil {
			return "", err
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func Login(username string, password string) (*VerifiedUser, error) {
	user, err := database.GetUser(username, password)
	if err != nil {
		return &VerifiedUser{}, fmt.Errorf("invalid credentials: %v", err)
	}
	token, err := CreateToken(user.Username)
	if err != nil {
		return &VerifiedUser{}, fmt.Errorf("token error: %v", err)
	}
	verifiedUser := VerifiedUser{user.Username, token}
	return &verifiedUser, nil
}

func Register(username string, password string) (*VerifiedUser, error) {
	user, err := database.CreateUser(username, password)
	if err != nil {
		return &VerifiedUser{}, fmt.Errorf("cannot register: %v", err)
	}
	token, err := CreateToken(user.Username)
	if err != nil {
		return &VerifiedUser{}, fmt.Errorf("token error: %v", err)
	}
	verifiedUser := VerifiedUser{user.Username, token}
	return &verifiedUser, nil
}
