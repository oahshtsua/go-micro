package main

import (
	"fmt"
	"net/http"
)

func (app *application) sendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
        return
	}

	msg := message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.sendSMTPMessage(msg)
    fmt.Println("[handler] Error sending mail?")
    fmt.Println(err)

	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
        return
	}

	payload := jsonResponse{
		Err: false,
		Msg: "Sent to " + requestPayload.To,
	}
	app.writeJSON(w, http.StatusAccepted, payload, nil)
}
