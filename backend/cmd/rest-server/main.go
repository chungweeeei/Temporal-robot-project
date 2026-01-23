package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chungweeeei/Temporal-robot-project/internal/api"
	config "github.com/chungweeeei/Temporal-robot-project/internal/config/api"
	"github.com/chungweeeei/Temporal-robot-project/internal/database"
	"go.temporal.io/sdk/client"
)

func main() {

	// Initialize database connection
	db := database.InitDB()

	// Initialize Temporal client
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}

	// Register restful server
	app := config.NewAppConfig(db, temporalClient)

	go listenForErrors(app)

	go listenForShutdown(app)

	router := api.NewRouter(app)

	app.InfoLog.Println("Starting REST server on :3000")
	if err := router.Run("localhost:3000"); err != nil {
		log.Panic(err)
	}
}

func listenForErrors(app *config.AppConfig) {

	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorDoneChan:
			return
		}
	}
}

func listenForShutdown(app *config.AppConfig) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Shutdown()
	os.Exit(0)
}
