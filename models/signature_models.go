package models

import (
	"database/sql"
	"time"
)

type Signatory struct {
	Name     string `json:"name" db:"name"`
	Position string `json:"position" db:"position"`
	Role     string `json:"role_sign" db:"role_sign"`
}

type Signatories struct {
	UUID             string         `json:"sign_uuid" db:"sign_uuid"`
	Name             string         `json:"name" db:"name"`
	Position         string         `json:"position" db:"position"`
	Role             string         `json:"role_sign" db:"role_sign"`
	IsSign           bool           `json:"is_sign" db:"is_sign"`
	IsApprove        sql.NullBool   `json:"is_approve" db:"is_approve"`
	ApprovalStatus   string         `json:"approval_status"`
	ReasonNotApprove string         `json:"reason_not_approve" db:"reason_not_approve"`
	Created_by       sql.NullString `json:"created_by" db:"created_by"`
	Created_at       time.Time      `json:"created_at" db:"created_at"`
	Updated_by       sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at       sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by       sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at       sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type Signatorie struct {
	UUID             string         `json:"sign_uuid" db:"sign_uuid"`
	Name             string         `json:"name" db:"name"`
	Position         string         `json:"position" db:"position"`
	Role             string         `json:"role_sign" db:"role_sign"`
	IsSign           bool           `json:"is_sign" db:"is_sign"`
	IsApprove        sql.NullBool   `json:"is_approve" db:"is_approve"`
	ReasonNotApprove string         `json:"reason_not_approve" db:"reason_not_approve"`
	Created_by       sql.NullString `json:"created_by" db:"created_by"`
	Created_at       time.Time      `json:"created_at" db:"created_at"`
	Updated_by       sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at       sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by       sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at       sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type UpdateSign struct {
	IsSign           bool      `json:"is_sign" db:"is_sign" validate:"required"`
	IsApprove        bool      `json:"is_approve" db:"is_approve"`
	ReasonNotApprove string    `json:"reason_not_approve" db:"reason_not_approve"`
	Updated_by       string    `json:"updated_by" db:"updated_by"`
	Updated_at       time.Time `json:"updated_at" db:"updated_at"`
}
