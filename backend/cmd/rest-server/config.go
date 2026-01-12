package main

import (
	"log"

	"go.temporal.io/sdk/client"
)

type Config struct {
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	ErrorChan      chan error
	ErrorDoneChan  chan bool
	TemporalClient client.Client
}
