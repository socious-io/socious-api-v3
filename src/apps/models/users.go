package models

import (
	"context"
	"socious/src/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID    `db:"id" json:"id"`
	FirstName           *string      `db:"first_name" json:"first_name"`
	LastName            *string      `db:"last_name" json:"last_name"`
	Username            string       `db:"username" json:"username"`
	Email               string       `db:"email" json:"email"`
	EmailText           *string      `db:"email_text" json:"email_text"`
	Phone               *string      `db:"phone" json:"phone"`
	WalletAddress       *string      `db:"wallet_address" json:"wallet_address"`
	Password            *string      `db:"password" json:"-"`
	RememberToken       *string      `db:"remember_token" json:"remember_token"`
	City                string       `db:"city" json:"city"`
	DescriptionSearch   *string      `db:"description_search" json:"description_search"`
	Address             *string      `db:"address" json:"address"`
	ExpiryDate          *time.Time   `db:"expiry_date" json:"expiry_date"`
	Status              string       `db:"status" json:"status"` // user_status as type, default 'INACTIVE'
	Mission             *string      `db:"mission" json:"mission"`
	Bio                 *string      `db:"bio" json:"-"`
	ViewAs              *int         `db:"view_as" json:"view_as"`
	JobTitle            *string      `db:"job_title" json:"job_title"`
	AvatarID            *uuid.UUID   `db:"avatar_id" json:"avatar_id"`
	PasswordExpired     bool         `db:"password_expired" json:"password_expired"`
	Language            *string      `db:"language" json:"language"`
	MyConversation      *string      `db:"my_conversation" json:"my_conversation"`
	ImpactPoints        float32      `db:"impact_points" json:"impact_points"`
	SocialCauses        *[]string    `db:"social_causes" json:"social_causes"` // social_causes_type[] as typ
	Followers           int          `db:"followers" json:"followers"`
	Followings          int          `db:"followings" json:"followings"`
	Avatar              *uuid.UUID   `db:"avatar" json:"avatar"`
	CoverImage          *uuid.UUID   `db:"cover_image" json:"cover_image"`
	Skills              []uint8      `db:"skills" json:"skills"`
	Country             string       `db:"country" json:"country"`
	MobileCountryCode   *string      `db:"mobile_country_code" json:"mobile_country_code"`
	OldId               *int         `db:"old_id" json:"old_id"`
	SearchTsv           string       `db:"search_tsv" json:"search_tsv"`
	Certificates        *[]string    `db:"certificates" json:"certificates"`
	Educations          *[]string    `db:"educations" json:"educations"`
	Goals               *string      `db:"goals" json:"goals"`
	GeonameId           int64        `db:"geoname_id" json:"geoname_id"`
	IsAdmin             bool         `db:"is_admin" json:"is_admin"`
	ProofspaceConnectId *string      `db:"proofspace_connect_id" json:"proofspace_connect_id"`
	OpenToWork          bool         `db:"open_to_work" json:"open_to_work"`
	OpenToVolunteer     bool         `db:"open_to_volunteer" json:"open_to_volunteer"`
	IdentityVerified    bool         `db:"identity_verified" json:"identity_verified"`
	IsContributor       *bool        `db:"is_contributor" json:"is_contributor"`
	Events              *[]uuid.UUID `db:"events" json:"events"`

	EmailVerifiedAt *time.Time `db:"email_verified_at" json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `db:"phone_verified_at" json:"phone_verified_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (User) FetchQuery() string {
	return "users/fetch"
}

func (u *User) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/register",
		u.FirstName, u.LastName, u.Username, u.Email, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Verify(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/verify",
		u.ID, u.Status,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) ExpirePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/expire_password",
		u.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdatePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_password",
		u.ID, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdateProfile(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_profile",
		u.ID, u.FirstName, u.LastName, u.Bio, u.JobTitle, u.Phone, u.Username,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func GetUser(id uuid.UUID) (*User, error) {
	u := new(User)
	if err := database.Fetch(u, id.String()); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_email", email); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByUsername(username string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_username", username); err != nil {
		return nil, err
	}
	return u, nil
}
