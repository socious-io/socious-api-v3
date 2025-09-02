package workers

import (
	"context"
	"log"
	"socious/src/apps/models"
)

func DeleteUser(form DeleteUserForm) error {
	ctx := context.Background()

	u, err := models.GetUser(form.User.ID)
	if err != nil {
		log.Printf("DeleteUser: Error fetching user: %v\n", err)
		return err
	}

	if err := u.Delete(ctx, form.Reason); err != nil {
		log.Printf("DeleteUser: Error deleting user: %v\n", err)
		return err
	}

	return nil
}

func SyncIdentities(form SyncForm) error {
	ctx := context.Background()

	user := models.GetTransformedUser(ctx, form.User)
	if err := user.Upsert(ctx); err != nil {
		return err
	}
	if err := user.AttachMedia(ctx, form.User); err != nil {
		log.Printf("SyncIdentities: Error attaching media to user: %v\n", err)
	}

	for _, o := range form.Organizations {
		organization := models.GetTransformedOrganization(ctx, o)
		if err := organization.Upsert(ctx, user.ID); err != nil {
			log.Printf("SyncIdentities: Error upserting org: %v\n", err)
			return err
		}
		if err := organization.AttachMedia(ctx, o, user.ID); err != nil {
			log.Printf("SyncIdentities: Error attaching media to org: %v\n", err)
			return err
		}
	}
	return nil
}
