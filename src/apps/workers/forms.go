package workers

import "github.com/socious-io/goaccount"

type SyncForm struct {
	Organizations []goaccount.Organization `json:"organizations"`
	User          goaccount.User           `json:"user" validate:"required"`
}

type DeleteUserForm struct {
	User   goaccount.User `json:"user" validate:"required"`
	Reason string         `json:"reason" validate:"required"`
}
