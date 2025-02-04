package data

import "database/sql"

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
