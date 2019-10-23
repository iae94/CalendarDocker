package api

import (
	cfg "calendar/pkg/config"
	"calendar/pkg/models"
	DB "calendar/pkg/psql"
	pb "calendar/pkg/services/api/gen"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	pt "github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	pg "github.com/jackc/pgtype"
	"go.uber.org/zap"
	"log"
)
// Server struct
type Server struct {
	Config   *cfg.APIConfig
	Logger   *zap.Logger
	Storage *DB.PSQLInterface
}



// Constructor
func NewAPIServer(config *cfg.APIConfig, logger *zap.Logger) (*Server, error) {
	dsn := fmt.Sprintf("user=%v password=%v dbname=%v host=%v sslmode=disable", config.API.DB.User, config.API.DB.Password, config.API.DB.DBName, config.API.DB.Host)
	//fmt.Printf("DSN:%v", dsn)
	dbInterface, err := DB.NewPSQLInterface(logger, dsn)
	if err != nil {
		log.Fatalf("Error during create PSQL instance: %v", err)
		return nil, err
	}

	return &Server{
		Config:      config,
		Logger:      logger,
		Storage: dbInterface,
	}, nil
}


// Methods to convert protobuf event to database event
func ToPBEvent(event *models.DBEvent) (*pb.Event, error) {
	var pbUUID string
	err := event.UUID.AssignTo(&pbUUID)
	if err != nil {
		log.Printf("ToPBEvent: Cannot scan UUID from DBEvent: %v", err)
		return nil, err
	}
	startDate, err := pt.TimestampProto(event.StartDate)
	if err != nil {
		log.Printf("ToPBEvent: Cannot convert StartDate from DBEvent: %v", err)
		return nil, err
	}
	endDate, err := pt.TimestampProto(event.EndDate)
	if err != nil {
		log.Printf("ToPBEvent: Cannot convert EndDate from DBEvent: %v", err)
		return nil, err
	}
	notifyTime, err := pt.TimestampProto(event.NotifyTime)
	if err != nil {
		log.Printf("ToPBEvent: Cannot convert NotifyTime from DBEvent: %v", err)
		return nil, err
	}

	return &pb.Event{
		UUID: pbUUID,
		Summary: event.Summary,
		Description: event.Description,
		User: event.User,
		StartDate: startDate,
		EndDate: endDate,
		NotifyTime: notifyTime,
	}, nil
}

func ToDBEvent(event *pb.Event) (*models.DBEvent, error) {
	var dbUUID pg.UUID

	err := dbUUID.Scan(event.UUID)
	if err != nil {
		log.Printf("FromPBEvent: Cannot scan UUID from PBEvent: %v", err)
		return nil, err
	}
	startDate, err := pt.Timestamp(event.StartDate)
	if err != nil {
		log.Printf("FromPBEvent: Cannot convert StartDate from PBEvent: %v", err)
		return nil, err
	}
	endDate, err := pt.Timestamp(event.EndDate)
	if err != nil {
		log.Printf("FromPBEvent: Cannot convert EndDate from PBEvent: %v", err)
		return nil, err
	}
	notifyTime, err := pt.Timestamp(event.NotifyTime)
	if err != nil {
		log.Printf("FromPBEvent: Cannot convert NotifyTime from PBEvent: %v", err)
		return nil, err
	}

	return &models.DBEvent{
		UUID: dbUUID,
		Summary: event.Summary,
		Description: event.Description,
		User: event.User,
		StartDate: startDate,
		EndDate: endDate,
		NotifyTime: notifyTime,
	}, nil
}




func (s *Server) EventCreate(ctx context.Context, event *pb.Event) (*pb.Response, error) {

	dbEvent, err := ToDBEvent(event)
	if err != nil {
		mess := fmt.Sprintf("Cannot convert PBEvent to DBEvent: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	s.Logger.Sugar().Infof("Try create event with uuid: %v", event.UUID)
	err = s.Storage.InsertEvent(dbEvent)
	if err != nil {
		mess := fmt.Sprintf("Storage: create event error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	return &pb.Response{Message: "Event created!"}, nil
}

func (s *Server) EventUpdate(ctx context.Context, request *pb.UpdateRequest) (*pb.Response, error) {

	dbEvent, err := ToDBEvent(request.Event)
	if err != nil {
		mess := fmt.Sprintf("Cannot convert PBEvent to DBEvent: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	s.Logger.Sugar().Infof("Try create event with uuid: %v", request.ID)
	err = s.Storage.UpdateEvent(request.ID, dbEvent)
	if err != nil {
		mess := fmt.Sprintf("Storage: update event error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}
	return &pb.Response{Message: "Event updated!"}, nil
}

func (s *Server) EventDelete(ctx context.Context, request *pb.DeleteRequest) (*pb.Response, error) {
	s.Logger.Sugar().Infof("Try delete event with uuid: %v", request.ID)
	err := s.Storage.DeleteEvent(request.ID)
	if err != nil {
		mess := fmt.Sprintf("Storage: delete event error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}
	return &pb.Response{Message: "Event deleted!"}, nil
}


func (s *Server) EventDayList(ctx context.Context, date *timestamp.Timestamp) (*pb.Events, error) {

	s.Logger.Info("Try get day events")
	timeDate, err := ptypes.Timestamp(date)
	if err != nil {
		mess := fmt.Sprintf("EventDayList converting time to protobuf time has failed: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}
	events, err := s.Storage.GetDayEvents(timeDate)
	if err != nil {
		mess := fmt.Sprintf("Storage: get day events error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	pbEvents := &pb.Events{Events: make([]*pb.Event, 0)}
	for _, ev := range events {
		pbEvent, err := ToPBEvent(ev)
		if err != nil {
			mess := fmt.Sprintf("Error during convert DBEvent to PBEvent: %v", err)
			s.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		pbEvents.Events = append(pbEvents.Events, pbEvent)
	}

	return pbEvents, nil
}

func (s *Server) EventWeekList(ctx context.Context, date *timestamp.Timestamp) (*pb.Events, error) {

	s.Logger.Info("Try get week events")
	timeDate, err := ptypes.Timestamp(date)
	if err != nil {
		mess := fmt.Sprintf("EventWeekList converting time to protobuf time has failed: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}
	events, err := s.Storage.GetWeekEvents(timeDate)
	if err != nil {
		mess := fmt.Sprintf("Storage: get week events error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	pbEvents := &pb.Events{Events: make([]*pb.Event, 0)}
	for _, ev := range events {
		pbEvent, err := ToPBEvent(ev)
		if err != nil {
			mess := fmt.Sprintf("Error during convert DBEvent to PBEvent: %v", err)
			s.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		pbEvents.Events = append(pbEvents.Events, pbEvent)
	}

	return pbEvents, nil

}

func (s *Server) EventMonthList(ctx context.Context, date *timestamp.Timestamp) (*pb.Events, error) {
	s.Logger.Info("Try get month events")
	timeDate, err := ptypes.Timestamp(date)
	if err != nil {
		mess := fmt.Sprintf("EventMonthList converting time to protobuf time has failed: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}
	events, err := s.Storage.GetMonthEvents(timeDate)
	if err != nil {
		mess := fmt.Sprintf("Storage: get month events error: %v", err)
		s.Logger.Info(mess)
		return nil, errors.New(mess)
	}

	pbEvents := &pb.Events{Events: make([]*pb.Event, 0)}
	for _, ev := range events {
		pbEvent, err := ToPBEvent(ev)
		if err != nil {
			mess := fmt.Sprintf("Error during convert DBEvent to PBEvent: %v", err)
			s.Logger.Sugar().Info(mess)
			return nil, errors.New(mess)
		}
		pbEvents.Events = append(pbEvents.Events, pbEvent)
	}

	return pbEvents, nil
}

