package database

import (
	"fmt"
)

type User struct {
	Username  string
	Password  string
	Watchlist []int
	Token     string
}

var database []User = []User{}

func GetUserByPassword(username string, password string) (*User, error) {
	for _, user := range database {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("username or password do not match")
}

func GetUserByToken(token string) (*User, error) {
	for _, user := range database {
		if user.Token == token {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("user token not found")
}

func GetUser(username string) (*User, error) {
	for _, user := range database {
		if user.Username == username {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("user not found")
}

func CreateUser(username string, password string) (*User, error) {
	for _, user := range database {
		if user.Username == username {
			return &User{}, fmt.Errorf("username already exists")
		}
	}
	newUser := User{Username: username, Password: password, Watchlist: []int{}}
	database = append(database, newUser)
	return &newUser, nil
}

func UpdateWatchlist(username string, watchlist []int) error {
	for i, item := range database {
		if item.Username == username {
			database[i].Watchlist = watchlist
			return nil
		}
	}
	return fmt.Errorf("could not update watchlist")
}

func AddToken(username string, token string) error {
	for i, user := range database {
		if user.Username == username {
			database[i].Token = token
			return nil
		}
	}
	return fmt.Errorf("error adding token")
}
