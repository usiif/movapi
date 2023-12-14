package data

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Movie MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movie: MovieModel{DB: db},
	}
}
