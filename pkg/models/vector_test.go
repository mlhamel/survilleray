package models

func (s *Suite) TestFindByPoint() {
	points, err := s.points.Find()

	s.NoError(err)

	s.NotEmpty(points)

	vectors, err := s.vectors.FindByPoint(&points[0])

	s.NoError(err)

	s.NotEmpty(vectors)
}
