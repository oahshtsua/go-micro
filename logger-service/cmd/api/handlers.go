package main

import (
	"fmt"
	"logger-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *application) writeLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload jsonPayload
	app.readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
    fmt.Println(event)
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
        return
	}

	resp := jsonResponse{
		Err: false,
		Msg: "entry successfully logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp, nil)
}
