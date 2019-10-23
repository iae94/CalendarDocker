package rabbit

import (
	cfg "calendar/pkg/config"
	"calendar/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Rabbit struct {
	Queue      amqp.Queue
	Channel    amqp.Channel
	Conn 	   amqp.Connection
	Storage    chan models.DBEvent
	Logger     *zap.Logger
	Config     *cfg.Rabbit
}

func NewRabbit(logger *zap.Logger, config *cfg.Rabbit) (*Rabbit, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/%v", config.User, config.Password, config.Host, config.Port, config.Vhost))
	if err != nil {
		logger.Sugar().Infof("amqp dial to rabbit give error: %v", err)
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Sugar().Infof("connect to channel give error: %v", err)
		return nil, err
	}
	queue, err := ch.QueueDeclare(
		"calendar", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		logger.Sugar().Infof("connect to channel give error: %v", err)
		return nil, err
	}
	return &Rabbit{
		Queue:      queue,
		Channel:    *ch,
		Conn: 		*conn,
		Storage:    make(chan models.DBEvent),
		Logger:     logger,
		Config:     config,
	}, nil
}

func (rb *Rabbit) Receive() error {
	response, err := rb.Channel.Consume(
		rb.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		rb.Logger.Sugar().Infof("Consuming error: %v", err)
		return err
	}

	forever := make(chan bool)

	go func() {
		for resp := range response {
			var event models.DBEvent
			err = json.Unmarshal(resp.Body, &event)
			if err != nil {
				rb.Logger.Sugar().Infof("Unmarshaling event error: %v", err)
			} else {
				rb.Logger.Sugar().Infof("Receive event: %v", event)
				rb.Storage <- event
			}
		}
	}()

	rb.Logger.Info("Consumer awaiting...")
	<-forever
	return nil
}

func (rb *Rabbit) Publish(event models.DBEvent) error {

	ev, err := json.Marshal(event)
	if err != nil {
		rb.Logger.Sugar().Infof("event marshaling give error: %v", err)
		return err
	}

	err = rb.Channel.Publish(
		"",
		rb.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        ev,
		})

	var pbUUID string
	_ = event.UUID.AssignTo(&pbUUID)
	rb.Logger.Sugar().Infof("Publish event with uuid: %v", pbUUID)
	if err != nil {
		rb.Logger.Sugar().Infof("Publish event give error: %v", err)
		return err
	}
	return nil
}

func (rb *Rabbit) Close() error {
	err := rb.Channel.Close()
	if err != nil {
		rb.Logger.Sugar().Infof("closing rabbit channel give error: %v", err)
		return err
	}
	err = rb.Conn.Close()
	if err != nil {
		rb.Logger.Sugar().Infof("closing rabbit connection give error: %v", err)
		return err
	}
	return nil
}




