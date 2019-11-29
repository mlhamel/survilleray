package models

func (s *Suite) TestFindByPoint() {
	pointRepos := NewPointRepository(s.context)
	vectorRepos := NewVectorRepository(s.context)

	points, err := pointRepos.FindByIcao24("c07c71")

	s.NoError(err)
	s.NotEmpty(points)

	vectors, err := vectorRepos.FindByPoint(&points[0])

	s.NoError(err)
	s.NotEmpty(vectors)
}
