package config

import (
	"database/sql"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
)

func NewTestDatabase(dsn string) (*gorm.DB, error) {
	txdb.Register("txdb", "postgres", dsn)

	db, err := sql.Open("txdb", "tx_1")

	if err != nil {
		return nil, err
	}

	orm, err := gorm.Open("postgres", db)

	if err != nil {
		return nil, err
	}

	orm.LogMode(true)

	return orm, nil
}

func NewDatabase(dsn string) (*gorm.DB, error) {
	orm, err := gorm.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	return orm, nil
}
