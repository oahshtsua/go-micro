package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   authPayload `json:"auth,omitempty"`
	Log    logPayload  `json:"log,omitempty"`
	Mail   mailPayload `json:"mail,omitempty"`
}

type authPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type logPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type mailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *application) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Err: false,
		Msg: "Hit the broker!",
	}
	err := app.writeJSON(w, http.StatusAccepted, payload, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *application) handleSubmit(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, http.StatusBadRequest, errors.New("unknown action"))
	}
}

func (app *application) authenticate(w http.ResponseWriter, a authPayload) {
	data, err := json.Marshal(a)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}

	request, err := http.NewRequest("POST", "http://authentication-service:8000/authenticate", bytes.NewBuffer(data))
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, http.StatusUnauthorized, errors.New("error calling auth service"))
		return
	}

	var jsonFromAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromAuth)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}

	if jsonFromAuth.Err {
		app.errorJSON(w, http.StatusUnauthorized, err)
	}

	var payload jsonResponse
	payload.Err = false
	payload.Msg = "Authenticated!"
	payload.Data = jsonFromAuth.Data

	app.writeJSON(w, http.StatusAccepted, payload, nil)
}

func (app *application) logItem(w http.ResponseWriter, entry logPayload) {
	js, err := json.Marshal(entry)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}

	request, err := http.NewRequest("POST", "http://logger-service:8000/log", bytes.NewBuffer(js))
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, response.StatusCode, errors.New("bad request"))
		return
	}

	var payload jsonResponse
	payload.Err = false
	payload.Msg = "Entry logged successfully!"

	app.writeJSON(w, http.StatusAccepted, payload, nil)
}

func (app *application) sendMail(w http.ResponseWriter, m mailPayload) {
	js, err := json.Marshal(m)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}

	request, err := http.NewRequest("POST", "http://mailer-service:8000/send", bytes.NewBuffer(js))
    fmt.Println(err)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, response.StatusCode, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Err = false
	payload.Msg = "sent email to: " + m.To

	app.writeJSON(w, http.StatusAccepted, payload, nil)
}
