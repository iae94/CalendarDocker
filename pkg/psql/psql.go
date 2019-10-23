package psql

import (
	"calendar/pkg/models"
	"context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
	"time"
)

type PSQLInterface struct {
	DB *sqlx.DB
	Logger *zap.Logger
}

func NewPSQLInterface(logger *zap.Logger, dsn string) (*PSQLInterface, error) {

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Printf("Cannot connect to DB: %v", err)
		return nil, err
	}
	return &PSQLInterface{
		DB: db,
		Logger: logger,
	}, nil
}


func (adb *PSQLInterface) InsertEvent(event *models.DBEvent) error{

	var UUID string
	err := event.UUID.AssignTo(&UUID)

	existingEvent, err := adb.GetEventByUUID(UUID)
	if err != nil {
		mess := fmt.Sprintf("Cannot check existence of event with UUID: %v", event.UUID)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}

	if existingEvent != nil {
		return errors.New("event with such UUID already exists")
	}

	insertSQL := "INSERT INTO public.\"Events\" (\"Summary\", \"Description\", \"User\", \"StartDate\", \"EndDate\", \"NotifyTime\", \"UUID\") VALUES (:Summary, :Description, :User, :StartDate, :EndDate, :NotifyTime, :UUID)"

	ctx := context.Background()
	_, err = adb.DB.NamedExecContext(ctx, insertSQL, event)
	if err != nil {
		mess := fmt.Sprintf("Cannot insert new event in to DB: %v", err)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}

	return nil
}


func (adb *PSQLInterface) DeleteEvent(eventUUID string) error {

	ctx := context.Background()

	deleteSQL := "DELETE FROM public.\"Events\" where \"UUID\" = :uuid"
	res, err := adb.DB.NamedExecContext(ctx, deleteSQL, map[string]interface{}{
		"uuid": eventUUID,
	})
	if err != nil {
		mess := fmt.Sprintf("Cannot delete event with: %v", err)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}

	deletedRows, err := res.RowsAffected()
	if err != nil {
		mess := fmt.Sprintf("Error during extract RowsAffected: %v", err)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}
	adb.Logger.Sugar().Infof("%v rows was deleted", deletedRows)
	return nil
}


func (adb *PSQLInterface) UpdateEvent(eventUUID string, newEvent *models.DBEvent) error {

	res, err := adb.DB.NamedExecContext(
		context.Background(),
		"UPDATE public.\"Events\" SET \"Summary\"=:summary, \"Description\"=:description, \"User\"=:user, \"StartDate\"=:startDate, \"EndDate\"=:endDate, \"NotifyTime\"=:notifyTime, \"UUID\"=:uuid WHERE \"UUID\"=:keyUUID",
		map[string]interface{}{
			"summary": newEvent.Summary,
			"description": newEvent.Description,
			"user": newEvent.User,
			"startDate": newEvent.StartDate,
			"endDate": newEvent.NotifyTime,
			"notifyTime": newEvent.NotifyTime,
			"uuid": newEvent.UUID,
			"keyUUID": eventUUID,
		},
	)

	if err != nil {
		mess := fmt.Sprintf("Cannot update event: %v", err)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}

	updatedRows, err := res.RowsAffected()
	if err != nil {
		mess := fmt.Sprintf("Error during extract RowsAffected: %v", err)
		adb.Logger.Info(mess)
		return errors.New(mess)
	}
	adb.Logger.Sugar().Infof("%v rows was updated", updatedRows)
	if updatedRows == 0 {
		return errors.New("event with such UUID does not exist in DB")
	}

	return nil
}


func (adb *PSQLInterface) GetEventByUUID(eventUUID string) (*models.DBEvent, error) {

	//var UUID pgtype.UUID
	//err := UUID.Set(eventUUID)

	ctx := context.Background()

	findSQL := "SELECT * FROM public.\"Events\" where \"UUID\" = :uuid"
	rows, err := adb.DB.NamedQueryContext(ctx, findSQL, map[string]interface{}{
		"uuid": eventUUID,
	})
	defer rows.Close()

	if err != nil {
		adb.Logger.Sugar().Infof("Error during sql select event by uuid: %v", err)
		return nil, err
	}

	var count = 0
	var scanEvent models.DBEvent
	for rows.Next() {
		scanEvent = models.DBEvent{} // nil each iteration
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			//adb.Logger.Sugar().Info(scanEvent)
			return nil, errors.New(mess)
		}
		count += 1
	}

	if count == 0 {
		adb.Logger.Sugar().Infof("No event with uuid: %v in DB", eventUUID)
		return nil, nil
	}

	return &scanEvent, nil
}


func (adb *PSQLInterface) GetEventsByUser(user string) ([]*models.DBEvent, error) {

	events := make([]*models.DBEvent, 0)


	ctx := context.Background()

	findSQL := "SELECT * FROM public.\"Events\" where \"User\" = :uuid"
	rows, err := adb.DB.NamedQueryContext(ctx, findSQL, map[string]interface{}{
		"uuid": user,
	})
	defer rows.Close()

	if err != nil {
		mess := fmt.Sprintf("Error during sql select event by user: %v", err)
		adb.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	var count = 0
	var scanEvent models.DBEvent
	for rows.Next() {
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		count += 1

		events = append(events, &scanEvent)
	}


	if count == 0 {
		adb.Logger.Sugar().Infof("No events for user: %v in DB", user)
		return nil, nil
	}

	return events, nil
}


func (adb *PSQLInterface) GetAllEvents() ([]*models.DBEvent, error) {

	events := make([]*models.DBEvent, 0)


	ctx := context.Background()

	findSQL := "SELECT * FROM public.\"Events\""
	rows, err := adb.DB.NamedQueryContext(ctx, findSQL, map[string]interface{}{})
	defer rows.Close()

	if err != nil {
		mess := fmt.Sprintf("Error during sql select event by user: %v", err)
		adb.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	var count = 0
	for rows.Next() {
		var scanEvent models.DBEvent
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		count += 1

		events = append(events, &scanEvent)
	}


	if count == 0 {
		adb.Logger.Info("No events in DB!")
		return nil, nil
	}

	return events, nil
}



func (adb *PSQLInterface) GetDayEvents(currentTime time.Time) ([]*models.DBEvent, error) {

	events := make([]*models.DBEvent, 0)

	rows, err := adb.DB.NamedQueryContext(
		context.Background(),
		"SELECT * from public.\"Events\" where EXTRACT(year from public.\"Events\".\"StartDate\") = :year and EXTRACT(month from public.\"Events\".\"StartDate\") = :month and EXTRACT(day from public.\"Events\".\"StartDate\") = :day",
		map[string]interface{}{
			"year": currentTime.Year(),
			"month": currentTime.Month(),
			"day": currentTime.Day(),
		},
	)
	defer rows.Close()
	if err != nil {
		mess := fmt.Sprintf("Error during sql select day events: %v", err)
		adb.Logger.Info(mess)
		return nil, errors.New(mess)
	}


	var count = 0
	for rows.Next() {
		var scanEvent models.DBEvent
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		count += 1

		events = append(events, &scanEvent)
	}


	if count == 0 {
		adb.Logger.Sugar().Infof("No events in day of: %v in DB", currentTime)
		return nil, nil
	}


	return events, nil
}


func (adb *PSQLInterface) GetWeekEvents(currentTime time.Time) ([]*models.DBEvent, error) {

	events := make([]*models.DBEvent, 0)
	endWeekTime := currentTime.AddDate(0, 0, 7)

	rows, err := adb.DB.NamedQueryContext(
		context.Background(),
		"SELECT * from public.\"Events\" where DATE_TRUNC('day', public.\"Events\".\"StartDate\") between :startWeek and :endWeek",
		map[string]interface{}{
			"startWeek": currentTime.Truncate(time.Hour * 24),
			"endWeek": endWeekTime.Truncate(time.Hour * 24),
		},
	)
	defer rows.Close()
	if err != nil {
		mess := fmt.Sprintf("Error during sql select week events: %v", err)
		adb.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	var count = 0
	for rows.Next() {
		var scanEvent models.DBEvent
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		count += 1

		events = append(events, &scanEvent)
	}


	if count == 0 {
		adb.Logger.Sugar().Infof("No events in day of: %v in DB", currentTime)
		return nil, nil
	}

	return events, nil
}


func (adb *PSQLInterface) GetMonthEvents(currentTime time.Time) ([]*models.DBEvent, error) {
	events := make([]*models.DBEvent, 0)

	rows, err := adb.DB.NamedQueryContext(
		context.Background(),
		"SELECT * from public.\"Events\" where EXTRACT(year from public.\"Events\".\"StartDate\") = :year and EXTRACT(month from public.\"Events\".\"StartDate\") = :month",
		map[string]interface{}{
			"year": currentTime.Year(),
			"month": currentTime.Month(),
		},
	)
	defer rows.Close()
	if err != nil {
		mess := fmt.Sprintf("Error during sql select month events: %v", err)
		adb.Logger.Info(mess)
		return nil, errors.New(mess)
	}


	var count = 0
	for rows.Next() {
		var scanEvent models.DBEvent
		if err := rows.StructScan(&scanEvent); err != nil {
			mess := fmt.Sprintf("Error during scan event from ResultSet: %v", err)
			adb.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		count += 1

		events = append(events, &scanEvent)
	}


	if count == 0 {
		adb.Logger.Sugar().Infof("No events in month of: %v in DB", currentTime)
		return nil, nil
	}


	return events, nil

}