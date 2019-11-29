package models

import "time"

func (s *Suite) TestString() {
	repos := NewPointRepository(s.context)

	points, err := repos.Find()

	s.NoError(err, "Cannot transform point to string")
	s.Equal("(c07c72, NDL322, 1568688174.000000)", points[0].String())
}

func (s *Suite) TestFindOverlaps() {
	pointRepos := NewPointRepository(s.context)
	districtRepos := NewDistrictRepository(s.context)

	points, err := pointRepos.Find()
	districts, err := districtRepos.Find()
	first := districts[0]
	overlaps, err := points[0].FindOverlaps(first)

	s.NoError(err, "Cannot find overlaps")
	s.True(overlaps, "Overlaps is %t", overlaps)
}

func (s *Suite) TestFindByVectorizedAt() {
	pointRepos := NewPointRepository(s.context)

	points, err := pointRepos.FindByVectorizedAt(nil)

	s.NoError(err, "Cannot points without any vectorized_at")
	s.NotEmpty(points, "Cannot finds Points without any vectorized_at")

	vectorizedAt, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	points, err = pointRepos.FindByVectorizedAt(&vectorizedAt)

	s.NoError(err, "Cannot points without any vectorized_at")
	s.NotEmptyf(points, "Cannot finds Points at %s", vectorizedAt.String())
}
