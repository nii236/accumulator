package accumulator

import (
	"accumulator/db"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	vrc "github.com/nii236/vrchat-go/client"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// refreshFriendCache in the database
func refreshFriendCache(IntegrationID int, updateBlob bool) error {
	integration, err := db.FindIntegrationG(null.Int64From(int64(IntegrationID)))
	if err != nil {
		return err
	}
	client, err := vrc.NewClient(vrc.ReleaseAPIURL, integration.AuthToken, integration.APIKey)
	if err != nil {
		return err
	}
	vrcResult, err := client.FriendList(false)
	if err != nil {
		return err
	}

	vrcResult2, err := client.FriendList(true)
	if err != nil {
		return err
	}
	vrcResult = append(vrcResult, vrcResult2...)
	for _, vrcFriend := range vrcResult {
		newBlobFilename := uuid.Must(uuid.NewV4()).String()
		if updateBlob {
			resp, err := http.Get(vrcFriend.CurrentAvatarThumbnailImageURL)
			if err != nil {
				fmt.Println(err)
				continue
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}

			blob := &db.Blob{
				FileName:      newBlobFilename,
				MimeType:      "image/jpg",
				FileSizeBytes: int64(len(b)),
				EXTENSION:     "jpg",
				File:          b,
			}
			err = blob.InsertG(boil.Infer())
			if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
				return err
			}

		}
		record := &db.Friend{
			IntegrationID:                 int64(IntegrationID),
			VrchatID:                      vrcFriend.ID,
			VrchatUsername:                vrcFriend.Username,
			VrchatDisplayName:             vrcFriend.DisplayName,
			VrchatAvatarImageURL:          vrcFriend.CurrentAvatarImageURL,
			VrchatAvatarThumbnailImageURL: vrcFriend.CurrentAvatarThumbnailImageURL,
			VrchatLocation:                vrcFriend.Location,
		}
		if updateBlob {
			savedBlob, err := db.Blobs(db.BlobWhere.FileName.EQ(newBlobFilename)).OneG()
			if err != nil {
				return err
			}
			record.AvatarBlobID = savedBlob.ID
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
			db.FriendColumns.AvatarBlobID:                  vrcFriend,
			db.FriendColumns.VrchatUsername:                vrcFriend.Username,
			db.FriendColumns.VrchatDisplayName:             vrcFriend.DisplayName,
			db.FriendColumns.VrchatAvatarImageURL:          vrcFriend.CurrentAvatarImageURL,
			db.FriendColumns.VrchatAvatarThumbnailImageURL: vrcFriend.CurrentAvatarThumbnailImageURL,
			db.FriendColumns.VrchatLocation:                vrcFriend.Location,
		}
		if updateBlob {
			updateMany[db.FriendColumns.AvatarBlobID] = newBlobFilename
		}
		_, err = existing.UpdateAllG(updateMany)
		if err != nil {
			return err
		}
	}
	return nil
}
