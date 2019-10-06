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
