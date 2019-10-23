package main

import (
	pb "calendar/pkg/services/api/gen"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	pt "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	newUUID, _ := uuid.NewV4()
	date, _ := pt.TimestampProto( time.Date(2019, 10, 02, 9, 0, 0, 651387237, time.UTC))

	date2start, _ := pt.TimestampProto( time.Date(2019, 10, 23, 19, 0, 0, 651387237, time.UTC))
	date2end, _ := pt.TimestampProto( time.Date(2019, 10, 23, 21, 0, 0, 651387237, time.UTC))
	date2notify, _ := pt.TimestampProto( time.Date(2019, 10, 21, 16, 0, 0, 651387237, time.UTC))

	date3start, _ := pt.TimestampProto( time.Date(2019, 10, 26, 13, 0, 0, 651387237, time.UTC))
	date3end, _ := pt.TimestampProto( time.Date(2019, 10, 26, 14, 0, 0, 651387237, time.UTC))
	date3notify, _ := pt.TimestampProto( time.Date(2019, 10, 21, 15, 0, 0, 651387237, time.UTC))

	date4start, _ := pt.TimestampProto( time.Date(2017, 9, 23, 19, 0, 0, 651387237, time.UTC))
	date4end, _ := pt.TimestampProto( time.Date(2017, 9, 23, 21, 0, 0, 651387237, time.UTC))
	date4notify, _ := pt.TimestampProto( time.Date(2017, 9, 21, 16, 0, 0, 651387237, time.UTC))



	event := &pb.Event{
		User: "Vasya",
		//UUID: "3ebcf64b-1402-47da-ab23-3656b3268135",
		UUID: newUUID.String(),
		Summary: "Event1",
		StartDate: date,
		EndDate: date,
		Description: "Test event",
		NotifyTime: date,
	}
	event2 := &pb.Event{
		User: "Vasya",
		//UUID: "3ebcf64b-1402-47da-ab23-3656b3268135",
		UUID: newUUID.String(),
		Summary: "Soon_event1",
		StartDate: date2start,
		EndDate: date2end,
		Description: "Test event1",
		NotifyTime: date2notify,
	}
	event3 := &pb.Event{
		User: "Vasya",
		//UUID: "3ebcf64b-1402-47da-ab23-3656b3268135",
		UUID: newUUID.String(),
		Summary: "Soon_event2",
		StartDate: date3start,
		EndDate: date3end,
		Description: "Test event2",
		NotifyTime: date3notify,
	}
	event4 := &pb.Event{
		User: "Vasya",
		//UUID: "3ebcf64b-1402-47da-ab23-3656b3268135",
		UUID: newUUID.String(),
		Summary: "Old event",
		StartDate: date4start,
		EndDate: date4end,
		Description: "Test event3",
		NotifyTime: date4notify,
	}



	//ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	//defer cancel()
	ctx := context.Background()

	fmt.Printf("\nCreate simple event\n")
	createResp, err := client.EventCreate(ctx, event)
	if err != nil {
		log.Printf("Error when calling EventCreate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", createResp)

	fmt.Printf("\nCreate soon1 event\n")
	createResp, err = client.EventCreate(ctx, event2)
	if err != nil {
		log.Printf("Error when calling EventCreate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", createResp)

	fmt.Printf("\nCreate soon2 event\n")
	createResp, err = client.EventCreate(ctx, event3)
	if err != nil {
		log.Printf("Error when calling EventCreate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", createResp)

	fmt.Printf("\nCreate old event\n")
	createResp, err = client.EventCreate(ctx, event4)
	if err != nil {
		log.Printf("Error when calling EventCreate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", createResp)

	fmt.Printf("\nCreate event again\n")
	createResp, err = client.EventCreate(ctx, event)
	if err != nil {
		log.Printf("Error when calling EventCreate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", createResp)


	fmt.Printf("\nGet daily events\n")
	dayResp, err := client.EventDayList(ctx, date)
	if err != nil {
		log.Printf("Error when calling EventDayList: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", dayResp)

	log.Printf("\nGet weekly events\n")
	weekResp, err := client.EventWeekList(ctx, date)
	if err != nil {
		log.Printf("Error when calling EventWeekList: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", weekResp)

	fmt.Printf("\nGet monthly events\n")
	monthResp, err := client.EventMonthList(ctx, date)
	if err != nil {
		log.Printf("Error when calling EventMonthList: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", monthResp)


	fmt.Printf("\nUpdate event\n")
	event.Description = "Updated description"
	updateResp, err := client.EventUpdate(ctx, &pb.UpdateRequest{ID: event.UUID, Event: event})
	if err != nil {
		log.Printf("Error when calling EventUpdate: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", updateResp)


	fmt.Printf("\nGet monthly events\n")
	monthResp, err = client.EventMonthList(ctx, date)
	if err != nil {
		log.Printf("Error when calling EventMonthList: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", monthResp)


	//fmt.Printf("\nDelete event\n")
	//deleteResp, err := client.EventDelete(ctx, &pb.DeleteRequest{ID: event.UUID})
	//if err != nil {
	//	log.Printf("Error when calling EventDelete: %v\n", err)
	//}
	//fmt.Printf("Response from server: %v\n", deleteResp)


	fmt.Printf("\nGet weekly events\n")
	weekResp, err = client.EventWeekList(ctx, date)
	if err != nil {
		log.Printf("Error when calling EventWeekList: %v\n", err)
	}
	fmt.Printf("Response from server: %v\n", weekResp)

}
