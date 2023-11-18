package repositories

import "database/sql"

type RealEstateRepository struct {
	db *sql.DB
}

func NewRealEstateRepository(db *sql.DB) *RealEstateRepository {
	return &RealEstateRepository{db: db}
}
