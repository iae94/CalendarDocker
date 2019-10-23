module Calendar/cmd/api/server

go 1.12

require (
	calendar/pkg/config v0.0.0
	calendar/pkg/logger v0.0.0
	calendar/pkg/models v0.0.0
	calendar/pkg/psql v0.0.0
	calendar/pkg/rabbit v0.0.0
	calendar/pkg/services/api/gen v0.0.0
	calendar/pkg/services/api/server v0.0.0

	github.com/jackc/pgtype v1.0.1 // indirect
	github.com/jackc/pgx v3.6.0+incompatible // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/spf13/viper v1.4.0 // indirect
	google.golang.org/grpc v1.24.0
)

replace calendar/pkg/config v0.0.0 => ../../../pkg/config

replace calendar/pkg/logger v0.0.0 => ../../../pkg/logger

replace calendar/pkg/services/api/gen v0.0.0 => ../../../pkg/services/api/gen

replace calendar/pkg/services/api/server v0.0.0 => ../../../pkg/services/api/server

replace calendar/pkg/models v0.0.0 => ../../../pkg/models

replace calendar/pkg/psql v0.0.0 => ../../../pkg/psql

replace calendar/pkg/rabbit v0.0.0 => ../../../pkg/rabbit
