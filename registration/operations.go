package registration

import (
	"errors"
	"fmt"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) error {
	// parse submitted form values
	username := r.FormValue("username")
	password := r.FormValue("password")

	// check if it's taken already
	if isUsernameTaken(username) {
		return fmt.Errorf("The username %s is already taken!", username)
	}

	// can't use empty string
	if username == "" {
		return errors.New("You must enter a username!")
	}

	// add user
	tempUserList = append(tempUserList, User{Username: username, Password: password})
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) error {
	// parse submitted form values
	username := r.FormValue("username")
	password := r.FormValue("password")

	if verifyUserCredentials(username, password) {
		return nil
	}
	return errors.New("Your login credentials are incorrect!")
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	return nil
}
