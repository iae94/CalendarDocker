package scheduler

import (
	"calendar/pkg/config"
	"calendar/pkg/models"
	DB "calendar/pkg/psql"
	"calendar/pkg/rabbit"
	"fmt"
	"go.uber.org/zap"
	"log"
	"time"
)

type Scheduler struct {
	Config *config.SchedulerConfig
	Rabbit *rabbit.Rabbit
	PSQL *DB.PSQLInterface
	Logger   *zap.Logger
}

func NewScheduler(config *config.SchedulerConfig, logger *zap.Logger) (*Scheduler, error) {
	rabbitInstance , err := rabbit.NewRabbit(logger, &config.Scheduler.Rabbit)
	if err != nil {
		logger.Sugar().Infof("Scheduler rabbit initialization give error: %v", err)
		return nil, err
	}
	dsn := fmt.Sprintf("user=%v password=%v dbname=%v host=%v sslmode=disable", config.Scheduler.DB.User, config.Scheduler.DB.Password, config.Scheduler.DB.DBName, config.Scheduler.DB.Host)
	dbInterface, err := DB.NewPSQLInterface(logger, dsn)
	if err != nil {
		log.Fatalf("Error during create PSQL instance: %v", err)
		return nil, err
	}
	return &Scheduler{
		PSQL: dbInterface,
		Rabbit: rabbitInstance,
		Config: config,
		Logger: logger,
	}, nil
}

func (s *Scheduler) CleanEvents() {
	for {
		s.Logger.Info("Try remove old events...")
		sleepTime := time.Hour * 24 * time.Duration(s.Config.Scheduler.Cleaner.CleanDelay)
		events, err := s.PSQL.GetAllEvents()
		if err != nil {
			s.Logger.Sugar().Errorf("Scheduler -> CleanEvents -> GetAllEvents give error: %v", err)
			time.Sleep(sleepTime)
			continue
		}

		oldEvents := make([]*models.DBEvent, 0)
		currentTime := time.Now()
		for _, ev := range events {
			diff := currentTime.Sub(ev.StartDate) // Difference between now and event start date
			difference := int(diff.Hours()/ 24 / 365)
			if difference >= 1 {
				oldEvents = append(oldEvents, ev)
			}
		}

		if len(oldEvents) == 0 {
			s.Logger.Info("No old events was found!")
			time.Sleep(sleepTime)
			continue
		}

		deleted := 0

		s.Logger.Sugar().Infof("Find %v old events -> try delete", len(oldEvents))
		for _, ev := range oldEvents {

			// Convert pg.UUID to string
			var pbUUID string
			err := ev.UUID.AssignTo(&pbUUID)
			if err != nil {
				log.Printf("ToPBEvent: Cannot scan UUID from DBEvent: %v", err)
				time.Sleep(sleepTime)
				continue
			}

			err = s.PSQL.DeleteEvent(pbUUID)
			if err != nil {
				s.Logger.Sugar().Warnf("Deleting event give error: %v", err)
			} else {
				deleted += 1
			}
		}
		s.Logger.Sugar().Infof("Events was deleted: %v", deleted)

		s.Logger.Info("Going to sleep!")
		time.Sleep(sleepTime)
	}
}

func (s *Scheduler) SendEvents() {
	sleepTime := time.Second * 30

	for {
		s.Logger.Info("Check near events...")
		currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC)
		events, err := s.PSQL.GetMonthEvents(currentTime)
		if err != nil {
			s.Logger.Sugar().Errorf("Scheduler -> SendEvents -> GetMonthEvents give error: %v", err)
			continue
		}
		for _, ev := range events {

			isEventNear := currentTime.After(ev.NotifyTime)
			isNowBeforeEventStart := currentTime.Before(ev.StartDate)

			var pbUUID string
			_ = ev.UUID.AssignTo(&pbUUID)

			if isEventNear && isNowBeforeEventStart{
				s.Logger.Sugar().Infof("Put to rabbit near event with uuid: %v", pbUUID)
				err = s.Rabbit.Publish(*ev)
				if err != nil {
					s.Logger.Sugar().Errorf("Putting event give error: %v", err)
				}
			}
		}

		time.Sleep(sleepTime)

	}
}

func (s *Scheduler) Start() {

	forever := make(chan bool)

	// Cleaner goroutine
	go s.CleanEvents()
	// Rabbit pusher goroutine
	go s.SendEvents()

	s.Logger.Info("Start scheduler successfully!")
	<-forever

}