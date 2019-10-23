package models

import (
	"fmt"
	"time"
)
import pg "github.com/jackc/pgtype"


// Event struct for db communication
type DBEvent struct {
	ID             		int  		`db:"id"`
	UUID 				pg.UUID 	`db:"UUID"`
	Summary 			string		`db:"Summary"`
	Description         string		`db:"Description"`
	User                string		`db:"User"`
	StartDate           time.Time	`db:"StartDate"`
	EndDate             time.Time	`db:"EndDate"`
	NotifyTime          time.Time	`db:"NotifyTime"`
}

// String interface
func (dbe *DBEvent) String() string {
	var uuidStr string
	dbe.UUID.AssignTo(&uuidStr)
	return fmt.Sprintf("\nEvent:\n\tUUID: %v\n\tSummary: %v\n\tDescription: %v\n\tUser: %v\n\tStartDate: %v\n\tEndDate: %v\n\tNotifyTime: %v\n",
		uuidStr, dbe.Summary, dbe.Description, dbe.User, dbe.StartDate, dbe.EndDate, dbe.NotifyTime,
	)
}