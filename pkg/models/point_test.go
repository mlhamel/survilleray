package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestString() {
	points, err := s.points.Find()

	s.NoError(err, "Cannot transform point to string")
	s.Equal("(c07c71, NDL321, 1568688174.000000)", points[0].String())
}

func (s *Suite) TestFindOverlaps() {
	points, err := s.points.Find()
	districts, err := s.districts.Find()
	overlaps, err := points[0].FindOverlaps(districts[0])

	s.NoError(err, "Cannot find overlaps")
	s.True(overlaps, "Overlaps is %t", overlaps)
}
