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
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	points    PointRepository
	districts DistrictRepository
}

type testPointRepository struct {
	point *Point
}

type testDistrictRepository struct {
	district *District
}

func NewTestPointRepository() PointRepository {
	return &testPointRepository{point: &Point{
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
	}}
}

func (t *testPointRepository) Find() ([]*Point, error) {
	return []*Point{t.point}, nil
}

func (t *testPointRepository) FindByName(name string) (*Point, error) {
	return t.point, nil
}

func NewTestDistrictRepository() DistrictRepository {
	district, _ := NewDistrictFromJson("villeray", "../../data/districts/villeray.geojson")

	return &testDistrictRepository{district: district}
}

func (t *testDistrictRepository) Find() ([]*District, error) {
	return []*District{t.district}, nil
}

func (t *testDistrictRepository) FindByName(name string) (*District, error) {
	return t.district, nil
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

	s.points = NewTestPointRepository()
	s.districts = NewTestDistrictRepository()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
