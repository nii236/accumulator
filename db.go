package accumulator

import (
	"accumulator/bindata"
	"accumulator/db"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	migrate_bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func randomAvatar() ([]byte, error) {
	resp, err := http.Get("https://i.pravatar.cc/300")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func newMigrateInstance(conn *sqlx.DB) (*migrate.Migrate, error) {
	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})
	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return nil, fmt.Errorf("bindata instance: %w", err)
	}
	dbDriver, err := sqlite3.WithInstance(conn.DB, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("db instance: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("migrate instance: %w", err)
	}
	return m, nil
}
func Migrate(conn *sqlx.DB) error {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
func Drop(conn *sqlx.DB) error {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = m.Drop()
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
func Version(conn *sqlx.DB) (uint, bool, error) {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return 0, false, fmt.Errorf("migrate: %w", err)
	}
	v, d, err := m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("migrate: %w", err)
	}
	return v, d, nil
}
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
			b, err := randomAvatar()
			if err != nil {
				return err
			}
			blobFilename := uuid.Must(uuid.NewV4()).String()
			blob := &db.Blob{
				FileName:      blobFilename,
				MimeType:      "image/jpg",
				FileSizeBytes: int64(len(b)),
				EXTENSION:     "jpg",
				File:          b,
			}
			err = blob.InsertG(boil.Infer())
			if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
				return err
			}

			friend := friendFactory(integration.ID.Int64)
			if i == 0 {
				friend.IsTeacher = true
			}
			savedBlob, err := db.Blobs(db.BlobWhere.FileName.EQ(blobFilename)).OneG()
			if err != nil {
				return err
			}
			friend.AvatarBlobID = savedBlob.ID

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
