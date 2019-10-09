package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestString() {
	vectors, err := s.vectors.Find()

	s.NoError(err, "Cannot transform vector to string")
	s.Equal("(c07c71, NDL321, 1568688174.000000)", vectors[0].String())
}

func (s *Suite) TestFindOverlaps() {
	vectors, err := s.vectors.Find()
	districts, err := s.districts.Find()
	overlaps, err := vectors[0].FindOverlaps(districts[0])

	s.NoError(err, "Cannot find overlaps")
	s.True(overlaps, "Overlaps is %t", overlaps)
}
