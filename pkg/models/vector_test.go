package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestString() {
	s.Equal("(c07c71, NDL321, 1568688174.000000)", s.vector.String())
}

func (s *Suite) TestOverlaps() {
	overlaps, err := s.vector.Overlaps(s.DB, s.district)

	s.NoError(err, "Cannot find overlaps")
	s.True(overlaps, "Overlaps is %t", overlaps)
}
