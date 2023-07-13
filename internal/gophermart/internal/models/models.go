package models

import "github.com/jackc/pgx/v5/pgtype"

type Order struct {
	ID          int64
	UserName    string
	Number      int64
	Status      string
	Accural     int
	UploadedAt  pgtype.Timestamptz
	ProcessedAt pgtype.Timestamptz
}
