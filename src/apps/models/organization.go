package models

import (
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Name              *string    `db:"name" json:"name"`
	Bio               *string    `db:"bio" json:"bio"`
	Description       *string    `db:"description" json:"description"`
	Email             *string    `db:"email" json:"email"`
	Phone             *string    `db:"phone" json:"phone"`
	City              *string    `db:"city" json:"city"`
	Type              string     `db:"type" json:"type"` //type -> organization_type DEFAULT 'OTHER'
	Address           *string    `db:"address" json:"address"`
	Website           *string    `db:"website" json:"website"`
	SocialCauses      []string   `db:"social_causes" json:"social_causes"` //type -> social_causes_type[]
	Followers         int        `db:"followers" json:"followers"`
	Followings        int        `db:"followings" json:"followings"`
	Country           *string    `db:"country" json:"country"`
	WalletAddress     *string    `db:"wallet_address" json:"wallet_address"`
	ImpactPoints      float64    `db:"impact_points" json:"impact_points"`
	Mission           *string    `db:"mission" json:"mission"`
	Culture           *string    `db:"culture" json:"culture"`
	Image             *uuid.UUID `db:"image" json:"image"`
	CoverImage        *uuid.UUID `db:"cover_image" json:"cover_image"`
	MobileCountryCode *string    `db:"mobile_country_code" json:"mobile_country_code"`
	CreatedBy         *uuid.UUID `db:"created_by" json:"created_by"`
	Shortname         string     `db:"shortname" json:"shortname"`
	// OldId *int `db:"old_id" json:"old_id"`
	Status string `db:"status" json:"status"` //type -> org_status DEFAULT 'ACTIVE'
	// SearchTsv tsvector `db:"search_tsv" json:"search_tsv"`
	OtherPartyId    *string `db:"other_party_id" json:"other_party_id"`
	OtherPartyTitle *string `db:"other_party_title" json:"other_party_title"`
	OtherPartyUrl   *string `db:"other_party_url" json:"other_party_url"`
	GeonameId       *int    `db:"geoname_id" json:"geoname_id"`
	VerifiedImpact  bool    `db:"verified_impact" json:"verified_impact"`
	Hiring          bool    `db:"hiring" json:"hiring"`
	Size            *string `db:"size" json:"size"`
	Industry        *string `db:"industry" json:"industry"`
	Did             *string `db:"did" json:"did"`
	Verified        bool    `db:"verified" json:"verified"`
	ImpactDetected  bool    `db:"impact_detected" json:"impact_detected"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (*Organization) Create() error {
	return nil
}

// func (*Organization) Update() error {
// 	id := u.ID
// 	if oauthSession != nil {
// 		//update profile
// 		err := oauthSession.UpdateUserProfile(u)
// 		if err != nil {
// 			return err
// 		}
// 		u.ID = id
// 	}
// }

func (*Organization) Remove() error {
	return nil
}

func (*Organization) UpdateDID() error {
	return nil
}

func (*Organization) ToggleHiring() error {
	return nil
}

func getAllOrganizations() ([]Organization, error) {
	result := []Organization{}
	return result, nil
}

func GetOrganization(id uuid.UUID, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func getManyOrganizations(ids []uuid.UUID, identity uuid.UUID) ([]Organization, error) {
	result := []Organization{}
	return result, nil
}

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func shortnameExistsOrganization(shortname string) (bool, error) {
	return false, nil
}

func searchOrganizations(query string) ([]Organization, error) { // Do we need to implement this?
	result := []Organization{}
	return result, nil
}
