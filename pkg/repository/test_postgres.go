package repository

import (
	"github.com/Saparly-go-dev/tagma"
	"gorm.io/gorm"
)

type TestPostgres struct {
	db *gorm.DB
}

// NewCityPostgres creates a new instance of CityPostgres with GORM DB
func NewTestPostgres(db *gorm.DB) *TestPostgres {
	return &TestPostgres{db: db}
}

// Create inserts a new city into the database and returns its ID
func (r *TestPostgres) Create(test tagma.Test) (int, error) {
	newtest := tagma.Test{
		Name: test.Name,
	}

	result := r.db.Create(&newtest)
	if result.Error != nil {
		return 0, result.Error
	}

	return newtest.Id, nil
}

// GetAll retrieves all cities from the database
func (r *TestPostgres) GetAll() ([]tagma.Test, error) {
	var tests []tagma.Test

	// Execute the query to find cities with id greater than 0
	result := r.db.Model(&tagma.Test{}).Where("id > ?", 0).Find(&tests)

	// Check for errors and log if necessary
	if result.Error != nil {
		// Optionally log the error here
		return nil, result.Error // Return nil for cities and the error
	}

	// Return the slice of cities and no error
	return tests, nil
}
