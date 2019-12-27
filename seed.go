package accumulator

import (
	"accumulator/db"
	"errors"
	"fmt"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func Seed() error {
	u := userFactory()
	u.Email = "jtnguyen236@gmail.com"
	u.PasswordHash = HashPassword("password")
	u.Role = roleAdmin
	err := u.InsertG(boil.Infer())
	if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
		return err
	}
	for i := 0; i < 2; i++ {
		u := userFactory()
		err = u.InsertG(boil.Infer())
		if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
			return err
		}
	}
	users, err := db.Users().AllG()
	if err != nil {
		return err
	}
	for _, user := range users {
		integration := integrationFactory(user.ID.Int64)
		err = integration.InsertG(boil.Infer())
		if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
			return err
		}
	}
	integations, err := db.Integrations().AllG()
	if err != nil {
		return err
	}
	for _, integration := range integations {
		for i := 0; i < 5; i++ {
			friend := friendFactory(integration.ID.Int64)
			if i == 0 {
				friend.IsTeacher = true
			}
			err = friend.InsertG(boil.Infer())
			if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
				return err
			}
		}
	}

	integations, err = db.Integrations().AllG()
	if err != nil {
		return err
	}
	for _, integration := range integations {
		friends, err := db.Friends(db.FriendWhere.IntegrationID.EQ(integration.ID.Int64)).AllG()
		if err != nil {
			return err
		}
		teacherID := null.NewInt64(0, false)
		for _, friend := range friends {
			if friend.IsTeacher {
				teacherID = friend.ID
			}
		}
		if !teacherID.Valid {
			return errors.New("no teacher in integration")
		}
		for _, friend := range friends {
			if friend.IsTeacher {
				continue
			}

			for i := 0; i < 10; i++ {
				attendance := attendanceFactory(integration.ID.Int64, friend.ID.Int64, teacherID.Int64)
				err = attendance.InsertG(boil.Infer())
				if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
					return err
				}
			}
		}
	}
	return nil
}
func integrationFactory(userID int64) *db.Integration {
	data := &db.Integration{
		UserID:    userID,
		Username:  faker.Email(),
		APIKey:    faker.UUIDDigit(),
		AuthToken: faker.UUIDDigit(),
	}

	return data
}
func userFactory() *db.User {

	data := &db.User{
		Email:        faker.Email(),
		PasswordHash: HashPassword("password"),
	}

	return data
}
func friendFactory(integrationID int64) *db.Friend {
	data := &db.Friend{
		IntegrationID:                 integrationID,
		IsTeacher:                     false,
		VrchatID:                      faker.Username(),
		VrchatUsername:                faker.Username(),
		VrchatDisplayName:             faker.Name(),
		VrchatAvatarImageURL:          faker.URL(),
		VrchatAvatarThumbnailImageURL: faker.URL(),
		VrchatLocation:                fmt.Sprintf("%f-%f", faker.Latitude(), faker.Longitude()),
	}

	return data
}
func attendanceFactory(
	integrationID int64,
	friendID int64,
	teacherID int64,
) *db.Attendance {
	data := &db.Attendance{
		Timestamp:     faker.RandomUnixTime(),
		IntegrationID: null.Int64From(integrationID),
		FriendID:      null.Int64From(friendID),
		TeacherID:     null.Int64From(teacherID),
		Location:      fmt.Sprintf("%f-%f", faker.Latitude(), faker.Longitude()),
	}

	return data
}
