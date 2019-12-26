package accumulator

import (
	"accumulator/db"
	"context"
	"fmt"
	"strings"
	"time"

	vrc "github.com/nii236/vrchat-go/client"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"
)

// RunAttendanceTracker starts the tracking service
func RunAttendanceTracker(ctx context.Context, stepMinutes int, log *zap.SugaredLogger) error {
	log.Infow("start attendance tracker")
	log.Info("running tracker")
	integrations, err := db.Integrations().AllG()
	if err != nil {
		return err
	}
	for _, integration := range integrations {
		if !integration.ID.Valid {
			fmt.Println("invalid integration ID, this should never happen")
			continue
		}
		trackAttendance(integration.ID.Int64, integration.AuthToken, integration.APIKey)
	}
	t := time.NewTicker(time.Duration(stepMinutes) * time.Minute)
	for {
		select {
		case <-t.C:
			log.Info("running tracker")
			for _, integration := range integrations {
				if !integration.ID.Valid {
					fmt.Println("invalid integration ID, this should never happen")
					continue
				}
				trackAttendance(integration.ID.Int64, integration.AuthToken, integration.APIKey)
			}
		}
	}
}

// trackAttendance in the database
func trackAttendance(integrationID int64, authToken, apiKey string) error {
	vrcClient, err := vrc.NewClient(vrc.ReleaseAPIURL, authToken, apiKey)
	if err != nil {
		return err
	}

	teachers, err := db.Friends(db.FriendWhere.IsTeacher.EQ(true)).AllG()
	if err != nil {
		return err
	}

	students, err := db.Friends(db.FriendWhere.IsTeacher.EQ(false)).AllG()
	if err != nil {
		return err
	}

	vrcfriends, err := vrcClient.FriendList(true)
	if err != nil {
		return err
	}

	for _, teacher := range teachers {
		currentLocation := ""
		for _, vrcfriend := range vrcfriends {
			if vrcfriend.ID == teacher.VrchatID {
				currentLocation = vrcfriend.Location
			}
		}
		if currentLocation := "" || currentLocation := "offline" {
			log.Errorw("could not get teacher location", "vrc_id", vrcfriend.ID)
			continue
		}
		for _, vrcfriend := range vrcfriends {
			for _, student := range students {
				if student.IntegrationID != integrationID {
					continue
				}
				if student.VrchatID != vrcfriend.ID {
					continue
				}
				if vrcfriend.Location == currentLocation {
					record := &db.Attendance{
						Timestamp:     time.Now().UnixNano(),
						IntegrationID: null.Int64From(integrationID),
						FriendID:      student.ID,
						TeacherID:     teacher.ID,
						Location:      currentLocation,
					}
					err = record.InsertG(boil.Infer())
					if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
						return err
					}
				}
			}
		}
	}

	return nil
}
