package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Err  bool   `json:"error"`
	Msg  string `json:"message"`
	Data any    `json:"data,omitempty"`
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, status int, err error) {
	var payload jsonResponse
	payload.Err = true
	payload.Msg = err.Error()

	app.writeJSON(w, status, payload, nil)
}
