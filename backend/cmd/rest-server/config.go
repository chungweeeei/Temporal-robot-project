package main

import (
	"log"

	"github.com/chungweeeei/Temporal-robot-project/cmd/rest-server/data"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

type Config struct {
	DB             *gorm.DB
	Model          data.Models
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	ErrorChan      chan error
	ErrorDoneChan  chan bool
	TemporalClient client.Client
}
