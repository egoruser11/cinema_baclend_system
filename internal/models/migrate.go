package models

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&PaymentMethod{},
		&Token{},
		&Movie{},
		&Genre{},
		&Premiere{},
		&Review{},
		&Order{},
	)
}

func SetupDatabase(db *gorm.DB) error {

	if err := Migrate(db); err != nil {
		return err
	}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedGenres(db); err != nil {
		return err
	}

	return nil
}
