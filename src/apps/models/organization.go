package models

import (
	"context"
	"socious/src/apps/utils"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/socious-io/goaccount"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID                uuid.UUID          `db:"id" json:"id"`
	Name              *string            `db:"name" json:"name"`
	Bio               *string            `db:"bio" json:"bio"`
	Description       *string            `db:"description" json:"description"`
	Email             *string            `db:"email" json:"email"`
	Phone             *string            `db:"phone" json:"phone"`
	City              *string            `db:"city" json:"city"`
	Type              string             `db:"type" json:"type"` //type -> organization_type DEFAULT 'OTHER'
	Address           *string            `db:"address" json:"address"`
	Website           *string            `db:"website" json:"website"`
	SocialCauses      pq.StringArray     `db:"social_causes" json:"social_causes"`
	Followers         int                `db:"followers" json:"followers"`
	Followings        int                `db:"followings" json:"followings"`
	Country           *string            `db:"country" json:"country"`
	WalletAddress     *string            `db:"wallet_address" json:"wallet_address"`
	ImpactPoints      float64            `db:"impact_points" json:"impact_points"`
	OldId             *string            `db:"old_id" json:"old_id"`
	SearchTSV         *string            `db:"search_tsv" json:"search_tsv"`
	Mission           *string            `db:"mission" json:"mission"`
	Culture           *string            `db:"culture" json:"culture"`
	MobileCountryCode *string            `db:"mobile_country_code" json:"mobile_country_code"`
	CreatedBy         *uuid.UUID         `db:"created_by" json:"created_by"`
	Shortname         string             `db:"shortname" json:"shortname"`
	Status            OrganizationStatus `db:"status" json:"status"`
	OtherPartyId      *string            `db:"other_party_id" json:"other_party_id"`
	OtherPartyTitle   *string            `db:"other_party_title" json:"other_party_title"`
	OtherPartyUrl     *string            `db:"other_party_url" json:"other_party_url"`
	GeonameId         *int               `db:"geoname_id" json:"geoname_id"`
	VerifiedImpact    bool               `db:"verified_impact" json:"verified_impact"`
	Hiring            bool               `db:"hiring" json:"hiring"`
	Size              *string            `db:"size" json:"size"`
	Industry          *string            `db:"industry" json:"industry"`
	Did               *string            `db:"did" json:"did"`
	Verified          bool               `db:"verified" json:"verified"`
	ImpactDetected    bool               `db:"impact_detected" json:"impact_detected"`

	LogoID   *uuid.UUID      `db:"logo_id" json:"logo_id"`
	Logo     *Media          `db:"-" json:"logo"`
	LogoJson *types.JSONText `db:"logo" json:"-"`
	Image    *uuid.UUID      `db:"image" json:"-"` //FIXME: temporary: we should unify it with other platforms

	CoverID    *uuid.UUID      `db:"cover_id" json:"cover_id"`
	Cover      *Media          `db:"-" json:"cover"`
	CoverJson  *types.JSONText `db:"cover" json:"-"`
	CoverImage *uuid.UUID      `db:"cover_image" json:"-"` //FIXME: temporary: we should unify it with other platforms

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func GetTransformedOrganization(ctx context.Context, org goaccount.Organization) *Organization {
	o := new(Organization)
	utils.Copy(org, o)

	if o.ID == uuid.Nil {
		newID, _ := uuid.NewUUID()
		o.ID = newID
	}

	if org.Verified || org.VerifiedImpact {
		o.Status = OrganizationStatusActive
	} else {
		o.Status = OrganizationStatusInactive
	}

	o.LogoID = nil
	o.CoverID = nil

	return o
}

func (o *Organization) AttachMedia(ctx context.Context, org goaccount.Organization, userID uuid.UUID) error {
	if org.Logo != nil {
		logo := new(Media)
		utils.Copy(org.Logo, logo)
		err := logo.Upsert(ctx)
		if err != nil {
			return err
		}
		o.LogoID = &logo.ID
	}

	if org.Cover != nil {
		cover := new(Media)
		utils.Copy(org.Cover, cover)
		err := cover.Upsert(ctx)
		if err != nil {
			return err
		}
		o.CoverID = &cover.ID
	}

	if err := o.Upsert(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (om *OrganizationMember) Create(ctx context.Context) error {
	_, err := database.Query(ctx, "organizations/create_member", om.OrganizationID, om.UserID)
	return err
}

func (o *Organization) Upsert(ctx context.Context, userID uuid.UUID) error {

	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}

	rows, err := database.TxQuery(
		ctx,
		tx,
		"organizations/upsert",
		o.ID,
		o.Shortname,
		o.Name,
		o.Bio,
		o.Description,
		o.Email,
		o.Phone,
		o.City,
		o.Country,
		o.Address,
		o.Website,
		o.Mission,
		o.Culture,
		o.Status,
		o.VerifiedImpact,
		o.Verified,
		o.LogoID,
		o.CoverID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(o); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := database.TxQuery(
		ctx,
		tx,
		"organizations/add_member",
		o.ID,
		userID,
	); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return database.Fetch(o, o.ID)
}

func GetOrganization(id uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func Member(orgID, userID uuid.UUID) (*OrganizationMember, error) {
	om := new(OrganizationMember)
	if err := database.Get(om, "organizations/get_member", orgID, userID); err != nil {
		return nil, err
	}
	return om, nil
}

func GetUserOrganizations(userId uuid.UUID) ([]Organization, error) {
	var (
		orgs      = []Organization{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("organizations/get_by_member", &fetchList, userId); err != nil {
		return orgs, err
	}

	if len(fetchList) < 1 {
		return orgs, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&orgs, ids...); err != nil {
		return orgs, err
	}
	return orgs, nil
}
