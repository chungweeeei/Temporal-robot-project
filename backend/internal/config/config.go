package config

import (
	"log"

	"github.com/chungweeeei/Temporal-robot-project/internal/database"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB             *gorm.DB
	Model          database.Models
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	ErrorChan      chan error
	ErrorDoneChan  chan bool
	TemporalClient client.Client
}

func NewAppConfig(db *gorm.DB, temporalClient client.Client) *AppConfig {
	infoLog := log.New(log.Writer(), "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(log.Writer(), "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &AppConfig{
		DB:             db,
		Model:          database.New(db),
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		ErrorChan:      make(chan error),
		ErrorDoneChan:  make(chan bool),
		TemporalClient: temporalClient,
	}
}

func (app *AppConfig) Shutdown() {

	// perform any cleanup tasks
	app.InfoLog.Println("Would run cleanup tasks...")

	// notify "listenForErrors" channel to close
	app.ErrorDoneChan <- true

	// shutdown
	app.InfoLog.Println("closing channels and shutting down application...")
	close(app.ErrorChan)
	close(app.ErrorDoneChan)
}
