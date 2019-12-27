package accumulator

import (
	"accumulator/db"
	"fmt"
	"strings"

	vrc "github.com/nii236/vrchat-go/client"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// refreshFriendCache in the database
func refreshFriendCache(IntegrationID int) error {
	integration, err := db.FindIntegrationG(null.Int64From(int64(IntegrationID)))
	if err != nil {
		return err
	}
	client, err := vrc.NewClient(vrc.ReleaseAPIURL, integration.AuthToken, integration.APIKey)
	if err != nil {
		return err
	}
	vrcResult, err := client.FriendList(true)
	if err != nil {
		return err
	}

	for _, vrcFriend := range vrcResult {
		// TODO: Handle blob
		// avatarBlob := &db.Blob{}
		// err = avatarBlob.InsertG()
		// if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
		// 	fmt.Println(err)

		// }

		record := &db.Friend{
			IntegrationID:                 int64(IntegrationID),
			VrchatID:                      vrcFriend.ID,
			VrchatUsername:                vrcFriend.Username,
			VrchatDisplayName:             vrcFriend.DisplayName,
			VrchatAvatarImageURL:          vrcFriend.CurrentAvatarImageURL,
			VrchatAvatarThumbnailImageURL: vrcFriend.CurrentAvatarThumbnailImageURL,
			VrchatLocation:                vrcFriend.Location,
		}
		existing, err := db.Friends(
			db.FriendWhere.VrchatID.EQ(vrcFriend.ID),
		).AllG()
		if len(existing) == 0 {
			err = record.InsertG(boil.Infer())
			if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
				fmt.Println(err)
				continue
			}
			continue
		}
		updateMany := db.M{
			db.FriendColumns.VrchatID:                      vrcFriend.ID,
			db.FriendColumns.VrchatUsername:                vrcFriend.Username,
			db.FriendColumns.VrchatDisplayName:             vrcFriend.DisplayName,
			db.FriendColumns.VrchatAvatarImageURL:          vrcFriend.CurrentAvatarImageURL,
			db.FriendColumns.VrchatAvatarThumbnailImageURL: vrcFriend.CurrentAvatarThumbnailImageURL,
			db.FriendColumns.VrchatLocation:                vrcFriend.Location,
		}
		_, err = existing.UpdateAllG(updateMany)
		if err != nil {
			return err
		}
	}
	return nil
}
