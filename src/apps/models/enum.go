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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
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
	ProjectLength1To3Day    ProjectLength = "1_3_DAYS"
	ProjectLength1Week      ProjectLength = "1_WEEK"
	ProjectLength2Weeks     ProjectLength = "2_WEEKS"
	ProjectLength1Month     ProjectLength = "1_MONTH"
)

func (pl *ProjectLength) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pl = ProjectLength(string(v))
	case string:
		*pl = ProjectLength(v)
	default:
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}

func (pl ProjectLength) Value() (driver.Value, error) {
	return string(pl), nil
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
		return fmt.Errorf("failed to scan type: %v", value)
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
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}

func (pk ProjectKind) Value() (driver.Value, error) {
	return string(pk), nil
}

type ContractStatus string

const (
	ContractStatusCreated          ContractStatus = "CREATED"
	ContractStatusClientApproved   ContractStatus = "CLIENT_APPROVED"
	ContractStatusSinged           ContractStatus = "SIGNED"
	ContractStatusProviderCanceled ContractStatus = "PROVIDER_CANCELED"
	ContractStatusClientCanceled   ContractStatus = "CLIENT_CANCELED"
)

func (pk *ContractStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pk = ContractStatus(string(v))
	case string:
		*pk = ContractStatus(v)
	default:
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}

type ContractType string

const (
	ContractTypeVolunteer ContractType = "VOLUNTEER"
	ContractTypePaid      ContractType = "PAID"
)

func (pk *ContractType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pk = ContractType(string(v))
	case string:
		*pk = ContractType(v)
	default:
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}

type ContractCommitmentPeriod string

const (
	ContractCommitmentHourly  ContractCommitmentPeriod = "HOURLY"
	ContractCommitmentDaily   ContractCommitmentPeriod = "DAILY"
	ContractCommitmentWeekly  ContractCommitmentPeriod = "WEEKLY"
	ContractCommitmentMonthly ContractCommitmentPeriod = "MONTHLY"
)

func (pk *ContractCommitmentPeriod) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pk = ContractCommitmentPeriod(string(v))
	case string:
		*pk = ContractCommitmentPeriod(v)
	default:
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}

type Currency string

const (
	USD Currency = "USD"
	JPY Currency = "JPY"
)

func (pk *Currency) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pk = Currency(string(v))
	case string:
		*pk = Currency(v)
	default:
		return fmt.Errorf("failed to scan type: %v", value)
	}
	return nil
}
