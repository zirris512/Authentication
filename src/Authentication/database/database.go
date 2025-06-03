package database

import (
	"fmt"
)

type User struct {
	Username string
	Password string
}

var database []User = []User{}

func GetUser(username string, password string) (*User, error) {
	for _, user := range database {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("username or password do not match")
}

func CreateUser(username string, password string) (*User, error) {
	for _, user := range database {
		if user.Username == username {
			return &User{}, fmt.Errorf("username already exists")
		}
	}
	newUser := User{Username: username, Password: password}
	database = append(database, newUser)
	return &newUser, nil
}
