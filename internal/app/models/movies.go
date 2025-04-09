package models

type MovieDTO struct {
	Name        string `db:"name"`
	Description *string `db:"description"`
	Path        *string `db:"path"`
}


type Movie struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	Path        string `db:"path"`
}

func FromDTO(m *MovieDTO) Movie {
	movie := Movie{}

	movie.Name = m.Name

	if m.Description != nil {
		movie.Description = *m.Description
	}
	if m.Path != nil {
		movie.Path = *m.Path
	}

	return movie
}