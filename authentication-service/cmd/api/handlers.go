package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
    fmt.Println("here in authenticate.")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		return
	}

	// validate user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, errors.New("Invalid credentials"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, http.StatusUnauthorized, errors.New("Invalid credentials"))
		return
	}

    fmt.Println("Reached here in authenticate. Time to log request.")
	// log authenticate
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in at %s", user.Email, time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")))
    fmt.Println(err)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, errors.New("Something went wrong"))
        return
	}

    fmt.Println("Request logged!")
	payload := jsonResponse{
		Err:  false,
		Msg:  fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}
	app.writeJSON(w, http.StatusAccepted, payload, nil)
}

func (app *application) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	js, err := json.Marshal(entry)
	if err != nil {
		return err
	}

    request, err := http.NewRequest("POST", "http://logger-service:8000/log", bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
