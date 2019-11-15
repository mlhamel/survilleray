package models

import (
	"database/sql/driver"
	"testing"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	cfg       *config.Config
	points    PointRepository
	vectors   VectorRepository
	districts DistrictRepository
	tx        driver.Tx
}

func TestVectorRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	orm, err := config.NewTestDatabase(config.GetEnv("DATABASE_URL", ""))

	if err != nil {
		panic(err)
	}

	s.cfg = config.NewConfigWithDB(orm)
	s.points = NewPointRepository(s.cfg)
	s.vectors = NewVectorRepository(s.cfg)
	s.districts = NewDistrictRepository(s.cfg)

	if err != nil {
		panic(err)
	}
}

func (s *Suite) SetupTest() {
	tx, err := s.cfg.DB().DB().Begin()

	if err != nil {
		panic(err)
	}

	s.tx = tx

	Migrate(s.cfg)

	point := Point{
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

	err = s.points.Insert(&point)

	district, err := NewDistrictFromJson("villeray", "villeray.geojson")

	if err != nil {
		panic(err)
	}

	err = s.districts.Insert(district)

	if err != nil {
		panic(err)
	}

	s.vectors.Insert(&Vector{
		Icao24:   "c07c71",
		CallSign: "NDL321",
		Country:  "Canada",
		Closed:   false,
		Points:   []Point{point},
	})

	if err != nil {
		panic(err)
	}
}

func (s *Suite) AfterTest(_, _ string) {
	s.tx.Rollback()
}
