package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type application struct {
	Mailer Mail
}

func main() {
	app := application{
		Mailer: createMail(),
	}
	srv := &http.Server{
		Addr:    ":8000",
		Handler: app.routes(),
	}
	log.Println("Starting mail service on port 8000")
	log.Fatal(srv.ListenAndServe())

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIN_PORT"))
	m := Mail{
		Domain:     os.Getenv("MAIL_DOMAIN"),
		Host:       os.Getenv("MAIL_HOST"),
		Port:       port,
		Username:   os.Getenv("MAIL_USERNAME"),
		Password:   os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName:   os.Getenv("MAIL_FROM_NAME"),
		FromAddr:   os.Getenv("MAIL_TO_ADDR"),
	}

	return m
}
