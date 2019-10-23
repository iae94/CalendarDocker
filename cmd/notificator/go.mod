module Calendar/cmd/notificator

go 1.12

require (
	calendar/pkg/config v0.0.0
	calendar/pkg/logger v0.0.0
	calendar/pkg/models v0.0.0
	calendar/pkg/rabbit v0.0.0
	calendar/pkg/services/notificator v0.0.0
	github.com/jackc/pgtype v1.0.1 // indirect
	github.com/spf13/viper v1.4.0 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271 // indirect
	go.uber.org/zap v1.11.0 // indirect
)

replace calendar/pkg/config v0.0.0 => ../../pkg/config

replace calendar/pkg/logger v0.0.0 => ../../pkg/logger

replace calendar/pkg/rabbit v0.0.0 => ../../pkg/rabbit

replace calendar/pkg/models v0.0.0 => ../../pkg/models

replace calendar/pkg/services/notificator v0.0.0 => ../../pkg/services/notificator
