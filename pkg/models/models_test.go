package models

import (
	"database/sql/driver"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/runtime"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	points    PointRepository
	vectors   VectorRepository
	districts DistrictRepository
	tx        driver.Tx
	context   *runtime.Context
}

func TestVectorRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	orm, err := config.NewTestDatabase(config.GetEnv("DATABASE_URL", ""))

	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig()
	s.context = runtime.NewContext(cfg, orm)

	s.points = NewPointRepository(s.context)
	s.vectors = NewVectorRepository(s.context)
	s.districts = NewDistrictRepository(s.context)

	if err != nil {
		panic(err)
	}
}

func (s *Suite) SetupTest() {
	tx, err := s.context.Database().DB().Begin()

	if err != nil {
		panic(err)
	}

	s.tx = tx

	Migrate(s.context)

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

	if err != nil {
		panic(err)
	}

	vector := Vector{
		Icao24:   "c07c71",
		CallSign: "NDL321",
		Country:  "Canada",
		Closed:   false,
	}

	s.vectors.Insert(&vector)

	if err != nil {
		panic(err)
	}

	err = s.vectors.AppendPoints(&vector, []*Point{&point})

	if err != nil {
		panic(err)
	}
}

func (s *Suite) AfterTest(_, _ string) {
	s.tx.Rollback()
}
