package main

import (
	"context"
	"fmt"
	handlers "handlersModule/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const (
	RETRIES_MAX = 5
)

func main() {
	cntRetries := 0
	dbHandler, err := handlers.SetupDB()
	if err != nil {
		for {
			dbHandler, err = handlers.SetupDB()
			if err != nil && cntRetries != RETRIES_MAX {
				cntRetries++
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
	}
	if err != nil {
		log.Fatal("Cannot set up db. Reason: ", err)
	}

	requestHandler := handlers.NewRequest(dbHandler)
	serveMux := mux.NewRouter()

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/encode", requestHandler.EncodeURL)
	postRouter.HandleFunc("/decode", requestHandler.DecodeURL)
	postRouter.Use(requestHandler.MiddlewareValidateData)

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{[a-z]+}", requestHandler.Redirect)

	serverLogger := log.New(os.Stdout, "ServerLog ", log.LstdFlags)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		serverLogger.Println("Starting server on port", os.Getenv("APP_PORT"))

		err := server.ListenAndServe()
		if err != nil {
			serverLogger.Printf("Error starting server: %s", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	serverLogger.Println("Received terminate, gracefull shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)

}
