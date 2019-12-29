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
func RunAttendanceTracker(ctx context.Context, d *Darer, stepMinutes int, log *zap.SugaredLogger) error {
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
		err := trackAttendance(d, integration.ID.Int64, integration.AuthToken, integration.AuthTokenNonce, integration.APIKey, log)
		if err != nil {
			log.Errorw(err.Error(), "integration_id", integration.ID.Int64, "integration_username", integration.Username)
			continue
		}
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
				err := trackAttendance(d, integration.ID.Int64, integration.AuthToken, integration.AuthTokenNonce, integration.APIKey, log)
				if err != nil {
					log.Errorw(err.Error(), "integration_id", integration.ID.Int64)
					continue
				}
			}
		}
	}
}

// trackAttendance in the database
func trackAttendance(d *Darer, integrationID int64, encryptedAuthToken []byte, nonce []byte, apiKey string, log *zap.SugaredLogger) error {
	err := refreshFriendCache(d, int(integrationID), false)
	if err != nil {
		return fmt.Errorf("could not refresh friend cache: %v", err)
	}

	decryptedAuthToken, err := d.decrypt(encryptedAuthToken, nonce)
	if err != nil {
		return err
	}
	vrcClient, err := vrc.NewClient(vrc.ReleaseAPIURL, string(decryptedAuthToken), apiKey)
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
		if currentLocation == "offline" {
			log.Errorw("teacher is offline", "vrc_id", teacher.VrchatID, "display_name", teacher.VrchatDisplayName)
			continue
		}
		if currentLocation == "" {
			log.Errorw("could not get teacher location", "vrc_id", teacher.VrchatID, "display_name", teacher.VrchatDisplayName)
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
						Timestamp:     time.Now().Unix(),
						IntegrationID: null.Int64From(integrationID),
						FriendID:      student.ID,
						TeacherID:     teacher.ID,
						Location:      currentLocation,
					}
					err = record.InsertG(boil.Infer())
					if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
						return fmt.Errorf("insert attendance record: %v", err)
					}
				}
			}
		}
	}

	return nil
}
