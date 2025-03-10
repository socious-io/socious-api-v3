package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/socious-io/goaccount"
	database "github.com/socious-io/pkg_database"
)

type OauthConnect struct {
	ID             uuid.UUID               `db:"id" json:"id"`
	IdentityId     uuid.UUID               `db:"identity_id" json:"identity_id"`
	Provider       OauthConnectedProviders `db:"provider" json:"provider"`
	MatrixUniqueId string                  `db:"matrix_unique_id" json:"matrix_unique_id"`
	AccessToken    string                  `db:"access_token" json:"access_token"`
	RefreshToken   *string                 `db:"refresh_token" json:"refresh_token"`
	Meta           *types.JSONText         `db:"meta" json:"meta"`
	Status         UserStatus              `db:"status" json:"status"`
	ExpiredAt      *time.Time              `db:"expired_at" json:"expired_at"`
	CreatedAt      time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at" json:"updated_at"`
}

func (OauthConnect) TableName() string {
	return "oauth_connects"
}

func (OauthConnect) FetchQuery() string {
	return "oauth_connects/fetch"
}

func GetOauthConnectByIdentityId(identityId uuid.UUID, provider OauthConnectedProviders) (*OauthConnect, error) {
	oc := new(OauthConnect)
	if err := database.Get(oc, "oauth_connects/get_by_identityid", identityId, provider); err != nil {
		return nil, err
	}
	return oc, nil
}

func GetOauthConnectByEmail(email string, provider OauthConnectedProviders) (*OauthConnect, error) {
	oc := new(OauthConnect)
	if err := database.Get(oc, "oauth_connects/get_by_email", email, provider); err != nil {
		return nil, err
	}
	return oc, nil
}

func (oc *OauthConnect) Create(ctx context.Context) error {
	rows, err := database.Query(ctx, "oauth_connects/create", oc.IdentityId, oc.Provider, oc.MatrixUniqueId, oc.AccessToken, oc.RefreshToken)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(oc); err != nil {
			return err
		}
	}
	return nil
}

func (oc *OauthConnect) Update(ctx context.Context) error {
	rows, err := database.Query(ctx, "oauth_connects/update", oc.ID, oc.MatrixUniqueId, oc.AccessToken, oc.RefreshToken)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(oc); err != nil {
			return err
		}
	}
	return nil
}

func (oc *OauthConnect) UpdateStatus(ctx context.Context, Status UserStatus) error {
	rows, err := database.Query(ctx, "oauth_connects/update_status", oc.ID, Status)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(oc); err != nil {
			return err
		}
	}
	return nil
}

func (oc *OauthConnect) SociousIdSession() goaccount.SessionToken {
	return goaccount.SessionToken{
		AccessToken:  oc.AccessToken,
		RefreshToken: *oc.RefreshToken,
		TokenType:    "Bearer",
	}
}
