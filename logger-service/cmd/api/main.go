package main

import (
	"context"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	rpcPort  = "5000"
	gRpcPort = "5001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type application struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println(err)
			panic(err)
		}
	}()

	app := application{
		Models: data.New(client),
	}
	// go app.serve()
	srv := &http.Server{
		Addr:    ":8000",
		Handler: app.routes(),
	}
	log.Println("Starting log service on port 8000...")
	log.Fatal(srv.ListenAndServe())
}

func (app *application) serve() {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: app.routes(),
	}
	log.Println("Starting log service on port 8000...")
	log.Fatal(srv.ListenAndServe())
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo: ", err)
		return nil, err
	}

    log.Println("Connected to mongo!")
	return conn, nil
}
