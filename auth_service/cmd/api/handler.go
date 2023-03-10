package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Authenticated accepts a json payload and attempts to authenticate a user
func (app *Config) Authenticated(w http.ResponseWriter, r *http.Request) {

	log.Println("Authenticated fun call")

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Data:    user,
		Error:   false,
	}

	log.Printf("Authenticated %v\n", payload)

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
