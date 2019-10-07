package models

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB       *gorm.DB
	mock     sqlmock.Sqlmock
	vector   *Vector
	district *District
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)
}

func (s *Suite) BeforeTest(_, _ string) {
	district, err := NewDistrictFromJson("villeray", "../../data/districts/villeray.geojson")

	s.NoError(err, "Error while creating test district")

	s.district = district
	s.vector = &Vector{
		Icao24:         "c07c71",
		CallSign:       "NDL321",
		OriginCountry:  "Canada",
		TimePosition:   1568688174,
		LastContact:    1568688174,
		Longitude:      -73.6275776,
		Latitude:       45.5339564,
		GeoAltitude:    0,
		OnGround:       true,
		Velocity:       9.26,
		Heading:        227,
		VerticalRate:   0,
		Sensors:        "EMPTY",
		BaroAltitude:   0,
		Squawk:         "3147",
		Spi:            false,
		PositionSource: 0,
	}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
