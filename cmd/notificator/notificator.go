package main

import (
	cfg "calendar/pkg/config"
	logging "calendar/pkg/logger"
	ntf "calendar/pkg/services/notificator"
	"log"
)

func main() {

	//Read config
	config, err := cfg.ReadNotificatorConfig()
	if err != nil {
		log.Fatalf("Reading notificator config error: %v \n", err)
	}

	logger, err := logging.CreateLogger(&config.Notificator.Logger)
	if err != nil {
		log.Fatalf("Creating notificator logger error: %v \n", err)
	}

	notificator, err := ntf.NewNotificator(config, logger)
	if err != nil {
		log.Fatalf("Cannot create notificator instance: %v", err)
	}

	notificator.Start()
}