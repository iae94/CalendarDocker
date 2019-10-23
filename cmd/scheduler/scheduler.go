package main

import (
	cfg "calendar/pkg/config"
	logging "calendar/pkg/logger"
	sch "calendar/pkg/services/scheduler"
	"log"
)

func main() {

	//Read config
	config, err := cfg.ReadSchedulerConfig()
	if err != nil {
		log.Fatalf("Reading scheduler config error: %v \n", err)
	}

	logger, err := logging.CreateLogger(&config.Scheduler.Logger)
	if err != nil {
		log.Fatalf("Creating scheduler logger error: %v \n", err)
	}

	scheduler, err := sch.NewScheduler(config, logger)
	if err != nil {
		log.Fatalf("Cannot create scheduler instance: %v", err)
	}

	scheduler.Start()
}