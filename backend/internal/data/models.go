package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict occured")
)

type Models struct {
	Containers ContainerModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Containers: ContainerModel{
			DB: db,
		},
	}
}
