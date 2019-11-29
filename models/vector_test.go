package models

func (s *Suite) TestFindByPoint() {
	pointRepos := NewPointRepository(s.context)
	vectorRepos := NewVectorRepository(s.context)

	points, err := pointRepos.Find()

	s.NoError(err)
	s.NotEmpty(points)

	vectors, err := vectorRepos.FindByPoint(&points[1])

	s.NoError(err)
	s.NotEmpty(vectors)
}
