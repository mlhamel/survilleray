package models

func (s *Suite) TestFindByName() {
	districts := NewDistrictRepository(s.cfg)
	district, err := districts.FindByName("villeray")

	s.NoError(err)
	s.NotNil(district)
	s.Equal(uint(1), district.ID)
	s.Empty(district.Points)
	s.Contains(district.Geometry, "MULTIPOLYGON(((-73.62195162 45.554561")
}
