package accumulator

import (
	"accumulator/db"
	"database/sql"
	"errors"

	vrc "github.com/nii236/vrchat-go/client"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// RefreshFriendCache in the database
func RefreshFriendCache(IntegrationID int) error {
	integration, err := db.FindIntegrationG(null.Int64From(int64(IntegrationID)))
	client, err := vrc.NewClient(vrc.ReleaseAPIURL, integration.AuthToken, integration.APIKey)
	if err != nil {
		return err
	}
	vrcResult, err := client.FriendList(true)
	if err != nil {
		return err
	}

	for _, vrcFriend := range vrcResult {
		record := &db.Friend{
			ID:                            vrcFriend.ID,
			IntegrationID:                 int64(IntegrationID),
			VrchatUsername:                vrcFriend.Username,
			VrchatDisplayName:             vrcFriend.DisplayName,
			VrchatAvatarImageURL:          vrcFriend.CurrentAvatarImageURL,
			VrchatAvatarThumbnailImageURL: vrcFriend.CurrentAvatarThumbnailImageURL,
		}
		_, err := db.FindFriendG(vrcFriend.ID)
		if errors.Is(err, sql.ErrNoRows) {
			record.InsertG(boil.Infer())
			continue
		}
		_, err = record.UpdateG(boil.Blacklist(db.FriendColumns.IsTeacher))
		if err != nil {
			return err
		}
	}
	return nil
}
