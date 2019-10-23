package notificator

import (
	"calendar/pkg/config"
	"calendar/pkg/rabbit"
	"go.uber.org/zap"
)

type Notificator struct {
	Config *config.NotificatorConfig
	Rabbit *rabbit.Rabbit
	Logger   *zap.Logger
}


func NewNotificator(config *config.NotificatorConfig, logger *zap.Logger) (*Notificator, error) {
	rabbitInstance , err := rabbit.NewRabbit(logger, &config.Notificator.Rabbit)
	if err != nil {
		logger.Sugar().Infof("Notificator rabbit initialization give error: %v", err)
		return nil, err
	}
	return &Notificator{
		Rabbit: rabbitInstance,
		Config: config,
		Logger: logger,
	}, nil
}



func (n *Notificator) SendWorker() {
	for event := range n.Rabbit.Storage {
		n.Logger.Sugar().Infof("\nБлижайшее событие у пользователя: %v\nЗапланировано в %v\nЗаканчивается в %v\nТема: %v\nОписание: %v\n\n", event.User, event.StartDate, event.EndDate, event.Summary, event.Description)
	}
}

func (n *Notificator) Listen() {
	err := n.Rabbit.Receive()
	if err != nil {
		n.Logger.Sugar().Errorf("Notificator rabbit consuming give error: %v", err)
	}
}

func (n *Notificator) Start() {


	forever := make(chan bool)

	// Rabbit receive listen goroutine
	go n.Listen()
	// Rabbit events channel listen goroutine
	go n.SendWorker()

	n.Logger.Info("Start notificator successfully!")
	<-forever
}
