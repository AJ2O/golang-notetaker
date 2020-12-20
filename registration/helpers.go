package registration

import (
	"fmt"
)

type User struct {
	Username string
	Password string
}

var tempUserList = []User{
	{Username: "user", Password: "pass"},
	{Username: "username", Password: "password"},
	{Username: "abc", Password: "123"},
}

func isUsernameTaken(userID string) bool {
	for curUserID := 0; curUserID < len(tempUserList); curUserID = curUserID + 1 {
		curUser := tempUserList[curUserID]
		if curUser.Username == userID {
			return true
		}
	}
	return false
}

func verifyUserCredentials(userID string, password string) bool {
	existingUser, err := getUser(userID)
	if err != nil {
		return false
	}
	return existingUser.Password == password
}

func getUser(userID string) (User, error) {
	for curUserID := 0; curUserID < len(tempUserList); curUserID = curUserID + 1 {
		curUser := tempUserList[curUserID]
		if curUser.Username == userID {
			return curUser, nil
		}
	}
	return User{}, fmt.Errorf("registration: no user exists with ID %s", userID)
}
