package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

type WorkSampleDocuments struct {
	Url      string `db:"url" json:"url"`
	Filename string `db:"filename" json:"filename"`
}

type WorkSampleType struct {
	ServiceID uuid.UUID `db:"service_id" json:"service_id"`
	Document  string    `db:"document" json:"document"`
}

type Project struct {
	ID                    uuid.UUID                `db:"id" json:"id"`
	IdentityID            uuid.UUID                `db:"identity_id" json:"identity_id"`
	Title                 *string                  `db:"title" json:"title"`
	Description           *string                  `db:"description" json:"description"`
	ProjectType           *ProjectType             `db:"project_type" json:"project_type"`
	ProjectLength         *ProjectLength           `db:"project_length" json:"project_length"`
	PaymentCurrency       *string                  `db:"payment_currency" json:"payment_currency"`
	PaymentRangeLower     *string                  `db:"payment_range_lower" json:"payment_range_lower"`
	PaymentRangeHigher    *string                  `db:"payment_range_higher" json:"payment_range_higher"`
	ExperienceLevel       *int                     `db:"experience_level" json:"experience_level"`
	Status                *ProjectStatus           `db:"status" json:"status"`
	PaymentType           *PaymentType             `db:"payment_type" json:"payment_type"`
	PaymentScheme         PaymentScheme            `db:"payment_scheme" json:"payment_scheme"`
	Country               *string                  `db:"country" json:"country"`
	Skills                pq.StringArray           `db:"skills" json:"skills"`
	CausesTags            pq.StringArray           `db:"causes_tags" json:"causes_tags"`
	OldId                 *int                     `db:"old_id" json:"old_id"`
	OtherPartyId          *string                  `db:"other_party_id" json:"other_party_id"`
	OtherPartyTitle       *string                  `db:"other_party_title" json:"other_party_title"`
	OtherPartyUrl         *string                  `db:"other_party_url" json:"other_party_url"`
	RemotePreference      *ProjectRemotePreference `db:"remote_preference" json:"remote_preference"`
	SearchTsv             string                   `db:"search_tsv" json:"search_tsv"`
	City                  *string                  `db:"city" json:"city"`
	WeeklyHoursLower      *string                  `db:"weekly_hours_lower" json:"weekly_hours_lower"`
	WeeklyHoursHigher     *string                  `db:"weekly_hours_higher" json:"weekly_hours_higher"`
	CommitmentHoursLower  *string                  `db:"commitment_hours_lower" json:"commitment_hours_lower"`
	CommitmentHoursHigher *string                  `db:"commitment_hours_higher" json:"commitment_hours_higher"`
	GeonameId             *int                     `db:"geoname_id" json:"geoname_id"`
	JobCategoryId         *uuid.UUID               `db:"job_category_id" json:"job_category_id"`
	ImpactJob             *bool                    `db:"impact_job" json:"impact_job"`
	Promoted              *bool                    `db:"promoted" json:"promoted"`
	ServiceTotalHours     *int                     `db:"service_total_hours" json:"service_total_hours"`
	ServicePrice          *int                     `db:"service_price" json:"service_price"`
	ServiceLength         *ServiceLength           `db:"service_length" json:"service_length"`
	Kind                  ProjectKind              `db:"kind" json:"kind"`
	WorkSamples           []WorkSampleDocuments    `db:"-" json:"work_samples"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	ExpiresAt *time.Time `db:"expires_at" json:"expires_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	WorkSamplesJson types.JSONText  `db:"work_samples" json:"-"`
	JobCategoryJson *types.JSONText `db:"job_category" json:"job_category"`
}

func (Project) TableName() string {
	return "projects"
}

func (Project) FetchQuery() string {
	return "projects/fetch"
}

func (p *Project) CreateService(ctx context.Context, workSamples []string) (*Project, error) {

	tx, err := database.GetDB().Beginx()
	rows, err := database.TxQuery(
		ctx,
		tx,
		"projects/create_service",
		p.IdentityID,
		p.Title,
		p.Description,
		p.PaymentCurrency,
		pq.Array(p.Skills),
		p.JobCategoryId,
		p.ServiceLength,
		p.ServiceTotalHours,
		p.ServicePrice,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(p); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rows.Close()

	workSamplesData := []WorkSampleType{}
	for _, workSample := range workSamples {
		workSamplesData = append(workSamplesData, WorkSampleType{ServiceID: p.ID, Document: workSample})
	}
	if len(workSamplesData) > 0 {
		if _, err = database.TxExecuteQuery(tx, "projects/create_work_samples", workSamplesData); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rows.Close()
	tx.Commit()

	s, err := GetService(p.ID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Project) UpdateService(ctx context.Context, workSamples []string) (*Project, error) {

	tx, err := database.GetDB().Beginx()
	rows, err := database.TxQuery(
		ctx,
		tx,
		"projects/update_service",
		p.ID,
		p.Title,
		p.Description,
		p.PaymentCurrency,
		pq.Array(p.Skills),
		p.JobCategoryId,
		p.ServiceLength,
		p.ServiceTotalHours,
		p.ServicePrice,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(p); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rows.Close()

	//delete and recreate files
	rows, err = database.TxQuery(ctx, tx, "projects/delete_work_samples",
		p.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	workSamplesData := []WorkSampleType{}
	for _, workSample := range workSamples {
		workSamplesData = append(workSamplesData, WorkSampleType{ServiceID: p.ID, Document: workSample})
	}
	if len(workSamplesData) > 0 {
		if _, err = database.TxExecuteQuery(tx, "projects/create_work_samples", workSamplesData); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rows.Close()
	tx.Commit()

	s, err := GetService(p.ID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Project) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "projects/delete", p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetServices(userId uuid.UUID, p database.Paginate) ([]Project, int, error) {
	var (
		services  = []Project{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("projects/get_services", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return services, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&services, ids...); err != nil {
		return nil, 0, err
	}
	return services, fetchList[0].TotalCount, nil
}

func GetService(id uuid.UUID) (*Project, error) {
	p := new(Project)
	if err := database.Fetch(p, id); err != nil {
		return nil, err
	}
	return p, nil
}
