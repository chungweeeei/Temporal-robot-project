package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
)

func main() {

	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		errorLog.Fatalf("Unable to create Temporal client: %v", err)
	}

	app := Config{
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		ErrorChan:      make(chan error),
		ErrorDoneChan:  make(chan bool),
		TemporalClient: temporalClient,
	}

	go app.listenForErrors()

	go app.listenForShutdown()

	app.serve()
}

func (app *Config) listenForErrors() {

	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorDoneChan:
			return
		}
	}
}

func (app *Config) listenForShutdown() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {

	// perform any cleanup tasks
	app.InfoLog.Println("Would run cleanup tasks...")

	// notify "listenForErrors" channel to close
	app.ErrorDoneChan <- true

	// shutdown
	app.InfoLog.Println("closing channels and shutting down application...")
	close(app.ErrorChan)
	close(app.ErrorDoneChan)
}

func (app *Config) serve() {

	srv := &http.Server{
		Addr:    "localhost:3000",
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting Shorten URL service")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
