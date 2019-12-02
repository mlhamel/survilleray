package models

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	points    PointRepository
	vectors   VectorRepository
	districts DistrictRepository
	tx        driver.Tx
	cfg       *config.Config
	ctx       context.Context
}

func TestVectorRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	database, err := config.NewTestDatabase(config.GetEnv("DATABASE_URL", ""))

	if err != nil {
		panic(err)
	}

	s.cfg = config.NewConfigWithDatabase(database)
}

func (s *Suite) SetupTest() {
	tx, err := s.cfg.Database().DB().Begin()

	if err != nil {
		panic(err)
	}

	s.tx = tx
	s.ctx = context.Background()

	err = Migrate(s.ctx, s.cfg)
	if err != nil {
		panic(err)
	}

	points, err := s.insertPoints()
	if err != nil {
		panic(err)
	}

	_, err = s.insertVectors(points)
	if err != nil {
		panic(err)
	}
}

func (s *Suite) insertVectors(points []Point) ([]Vector, error) {
	repos := NewVectorRepository(s.cfg)

	vector := Vector{
		Icao24:   "c07c71",
		CallSign: "NDL321",
		Country:  "Canada",
		Closed:   false,
	}

	err := repos.Insert(&vector)
	if err != nil {
		return []Vector{}, err
	}

	err = repos.AppendPoints(&vector, []Point{points[1]})
	if err != nil {
		return []Vector{}, err
	}

	return []Vector{vector}, nil
}

func (s *Suite) insertPoints() ([]Point, error) {
	repos := NewPointRepository(s.cfg)

	simplePoint := Point{
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

	err := repos.Insert(&simplePoint)
	if err != nil {
		return []Point{}, err
	}

	vectorizedAt, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	vectorizedPoint := Point{
		Icao24:         "c07c72",
		CallSign:       "NDL322",
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
		VectorizedAt:   vectorizedAt,
	}

	err = repos.Insert(&vectorizedPoint)
	if err != nil {
		return []Point{}, err
	}

	return []Point{simplePoint, vectorizedPoint}, nil
}

func (s *Suite) AfterTest(_, _ string) {
	s.tx.Rollback()
}
