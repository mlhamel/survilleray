package config

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
)

type database struct {
	orm *gorm.DB
}

func NewTestDatabase(dsn string) (*database, error) {
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

	return &database{orm}, nil
}

func NewDatabase(dsn string) (*database, error) {
	orm, err := newGoOrm(dsn)

	if err != nil {
		return nil, err
	}

	return &database{orm}, nil
}

func newGoOrm(dsn string) (*gorm.DB, error) {
	orm, err := gorm.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	orm.DB().SetConnMaxLifetime(time.Minute * 10)
	orm.DB().SetMaxIdleConns(10)
	orm.DB().SetMaxOpenConns(100)

	return orm, nil
}

func (d *database) DB() *gorm.DB {
	return d.orm
}

func (d *database) Close() error {
	if d.isDone() {
		return nil
	}

	d.orm.Lock()
	defer d.orm.Unlock()

	if err := d.orm.Close(); err != nil {
		return err
	}
	d.orm = nil
	return nil
}

func (d *database) isDone() bool {
	return d.orm != nil
}
