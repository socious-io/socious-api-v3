package models

import (
	"database/sql/driver"
	"fmt"
)

/*
Project
*/
type ProjectStatus string

const (
	ProjectStatusDraft  ProjectStatus = "DRAFT"
	ProjectStatusExpire ProjectStatus = "EXPIRE"
	ProjectStatusActive ProjectStatus = "ACTIVE"
)

func (ps *ProjectStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ps = ProjectStatus(string(v))
	case string:
		*ps = ProjectStatus(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (ps ProjectStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

type ProjectType string

const (
	ProjectTypeOneOff   ProjectType = "ONE_OFF"
	ProjectTypePartTime ProjectType = "PART_TIME"
	ProjectTypeFullTime ProjectType = "FULL_TIME"
)

func (pt *ProjectType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pt = ProjectType(string(v))
	case string:
		*pt = ProjectType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pt ProjectType) Value() (driver.Value, error) {
	return string(pt), nil
}

type PaymentModeType string

const (
	PaymentModeTypeCrypto PaymentModeType = "CRYPTO"
	PaymentModeTypeFiat   PaymentModeType = "FIAT"
)

func (pmt *PaymentModeType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pmt = PaymentModeType(string(v))
	case string:
		*pmt = PaymentModeType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pmt PaymentModeType) Value() (driver.Value, error) {
	return string(pmt), nil
}

type PaymentScheme string

const (
	PaymentSchemeHourly PaymentScheme = "HOURLY"
	PaymentSchemeFixed  PaymentScheme = "FIXED"
)

func (ps *PaymentScheme) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ps = PaymentScheme(string(v))
	case string:
		*ps = PaymentScheme(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (ps PaymentScheme) Value() (driver.Value, error) {
	return string(ps), nil
}

type PaymentService string

const (
	PaymentServiceStripe PaymentService = "STRIPE"
	PaymentServiceCrypto PaymentService = "CRYPTO"
)

func (ps *PaymentService) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ps = PaymentService(string(v))
	case string:
		*ps = PaymentService(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (ps PaymentService) Value() (driver.Value, error) {
	return string(ps), nil
}

type PaymentSourceType string

const (
	PaymentSourceTypeCard         PaymentSourceType = "CARD"
	PaymentSourceTypeCryptoWallet PaymentSourceType = "CRYPTO_WALLET"
)

func (pst *PaymentSourceType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pst = PaymentSourceType(string(v))
	case string:
		*pst = PaymentSourceType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pst PaymentSourceType) Value() (driver.Value, error) {
	return string(pst), nil
}

type PaymentType string

const (
	PaymentTypeVolunteer PaymentType = "VOLUNTEER"
	PaymentTypePaid      PaymentType = "PAID"
)

func (pt *PaymentType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pt = PaymentType(string(v))
	case string:
		*pt = PaymentType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pt PaymentType) Value() (driver.Value, error) {
	return string(pt), nil
}

type ProjectLength string

const (
	ProjectLengthLess1Day   ProjectLength = "LESS_THAN_A_DAY"
	ProjectLengthLess1Month ProjectLength = "LESS_THAN_A_MONTH"
	ProjectLength1To3Month  ProjectLength = "1_3_MONTHS"
	ProjectLength3To6Month  ProjectLength = "3_6_MONTHS"
	ProjectLengthMore6Month ProjectLength = "6_MONTHS_OR_MORE"
)

func (pl *ProjectLength) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pl = ProjectLength(string(v))
	case string:
		*pl = ProjectLength(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pl ProjectLength) Value() (driver.Value, error) {
	return string(pl), nil
}

type ServiceLength string

const (
	ServiceLengthLess1Day   ServiceLength = "LESS_THAN_A_DAY"
	ServiceLength1To3Day    ServiceLength = "1_3_DAYS"
	ServiceLength1Week      ServiceLength = "1_WEEK"
	ServiceLength2Weeks     ServiceLength = "2_WEEKS"
	ServiceLength1Month     ServiceLength = "1_MONTH"
	ServiceLength1To3Month  ServiceLength = "1_3_MONTHS"
	ServiceLength3To6Month  ServiceLength = "3_6_MONTHS"
	ServiceLengthMore6Month ServiceLength = "6_MONTHS_OR_MORE"
)

func (sl *ServiceLength) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*sl = ServiceLength(string(v))
	case string:
		*sl = ServiceLength(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (sl ServiceLength) Value() (driver.Value, error) {
	return string(sl), nil
}

type ProjectRemotePreference string

const (
	ProjectRemotePreferenceOnsite ProjectRemotePreference = "ONSITE"
	ProjectRemotePreferenceRemote ProjectRemotePreference = "REMOTE"
	ProjectRemotePreferenceHybrid ProjectRemotePreference = "HYBRID"
)

func (prp *ProjectRemotePreference) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*prp = ProjectRemotePreference(string(v))
	case string:
		*prp = ProjectRemotePreference(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (prp ProjectRemotePreference) Value() (driver.Value, error) {
	return string(prp), nil
}

type ProjectKind string

const (
	ProjectKindJob     ProjectKind = "JOB"
	ProjectKindService ProjectKind = "SERVICE"
)

func (pk *ProjectKind) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pk = ProjectKind(string(v))
	case string:
		*pk = ProjectKind(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (pk ProjectKind) Value() (driver.Value, error) {
	return string(pk), nil
}
