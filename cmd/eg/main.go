package main

import (
	example "expamle"
	"expamle/cmd/eg/services"
	"expamle/server"
	"fmt"
	kitLog "github.com/go-kit/kit/log"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	now := time.Now().UTC()

	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 35, 0, 0, time.UTC).Unix()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 30, 0, 0, time.UTC).Unix()

	fmt.Println("date:", base, "start:", start, "end:", end)

	// initializing configuration from environment variables
	config, err := example.NewConfiguration("example_session_string")
	if err != nil {
		log.Fatal(err)
		return
	}

	// connecting to postgres server
	db, err := gorm.Open(pg.Open(config.DNS), &gorm.Config{SkipDefaultTransaction: false})
	if err != nil {
		log.Fatal(err)
		return
	}

	// configuring logger
	var logger kitLog.Logger

	logger = kitLog.NewLogfmtLogger(kitLog.NewSyncWriter(os.Stderr))
	logger = kitLog.With(logger, "ts", kitLog.DefaultTimestampUTC)

	// AccountService
	account := services.InitAccountService(db, logger, config)

	//////////////////////////////////////////////
	// CreateMessage New Service With Given Components //
	//////////////////////////////////////////////
	// Service Logger
	sl := kitLog.With(logger, "component", "http")

	// Server instance
	svr := server.NewServer(sl, config)

	svr.InitAccountHandlers(account, config)

	// Start listening for http requests
	svr.Listen()

}
