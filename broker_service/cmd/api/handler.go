package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth   AuthPayload
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleSubmit handles a JSON payload that describes an action to take,
// processes it, and sends it where it needs to go
func (app *Config) HandleSubmit(w http.ResponseWriter, r *http.Request) {

	read, _ := io.ReadAll(r.Body)

	var requestPayload RequestPayload

	err := json.Unmarshal(read, &requestPayload)
	err = json.Unmarshal(read, &requestPayload.Auth)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		_ = app.errorJSON(w, errors.New("unknown action"))
	}

}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	//First Convert Golang Object(a) to []byte
	jsonData, _ := json.Marshal(a)

	//log.Println("kkt", jsonData)

	//Convert []bye value to  *bytes.Buffer to use in http request
	responseBody := bytes.NewBuffer(jsonData)

	//log.Printf("Value of response body %v", responseBody)
	authServiceURL := fmt.Sprintf("http://%s/authenticated", "auth_service")
	//Customizing the Request
	req, err := http.NewRequest("POST", authServiceURL, responseBody)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	//To send the request,  use the http.DefaultClient.DO
	res, err := http.DefaultClient.Do(req)

	log.Printf("response from auth %v\n", res.Body)

	if err != nil {

		_ = app.errorJSON(w, err, http.StatusBadRequest)

		return
	}

	defer res.Body.Close()

	// make sure we get back the right status code
	if res.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return

	} else if res.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	// create variable we'll read the response.Body from the authentication-service into
	var jsonFromService jsonResponse

	// decode the json we get from the authentication-service into our variable
	err = json.NewDecoder(res.Body).Decode(&jsonFromService)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//fmt.Printf("jsonFromService Message %v\n, Data is %v\n", jsonFromService.Message, jsonFromService.Data)

	//send json back to our end user, with user info embedded
	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
