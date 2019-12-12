package movie

import (
	"database/sql"
	"time"
)

type MovieService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *MovieService {
	return &MovieService{
		db: db,
	}
}

const (
	insertMovieQuery = "INSERT INTO movies (movie_title, movie_genre, movie_rating) VALUES (?, ?, ?)"
	selectMovieQuery = "SELECT id, movie_title, movie_genre, movie_rating, release_date FROM movies"
)

type Preferences struct {
	Genre          string
	AgeAppropriate int
	Date           time.Time
}

type Movie struct {
	Title       string
	Genre       string
	MinAge      int
	ReleaseDate time.Time
}

var SelectGenre = []string{
	"Sci-fi/Fantasy",
	"Animated",
	"Christmas",
	"Comedy",
	"Biography",
}

var movies = []Movie{
	Movie{
		Title:       "Star Wars: The Rise of Skywalker",
		Genre:       "Sci-fi/Fantasy",
		MinAge:      13,
		ReleaseDate: time.Date(2019, 12, 20, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Frozen II",
		Genre:       "Animated",
		MinAge:      0,
		ReleaseDate: time.Date(2019, 11, 22, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Last Christmas",
		Genre:       "Christmas",
		MinAge:      13,
		ReleaseDate: time.Date(2019, 11, 8, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Jumanji: The Next Level",
		Genre:       "Comedy",
		MinAge:      18,
		ReleaseDate: time.Date(2019, 12, 20, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "The Irishman",
		Genre:       "Biography",
		MinAge:      18,
		ReleaseDate: time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "A Beautiful day in the Neighborhood",
		Genre:       "Biography",
		MinAge:      13,
		ReleaseDate: time.Date(2019, 11, 22, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Ford v Ferrari",
		Genre:       "Action",
		MinAge:      13,
		ReleaseDate: time.Date(2019, 11, 15, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Charlie's Angels",
		Genre:       "Action",
		MinAge:      13,
		ReleaseDate: time.Date(2019, 11, 15, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Noelle",
		Genre:       "Christmas",
		MinAge:      0,
		ReleaseDate: time.Date(2019, 11, 12, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Terminator: Dark Fate",
		Genre:       "Sci-fi/Fantasy",
		MinAge:      18,
		ReleaseDate: time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Lady and the Tramp",
		Genre:       "Family",
		MinAge:      1,
		ReleaseDate: time.Date(2019, 11, 12, 0, 0, 0, 0, time.UTC),
	},
	Movie{
		Title:       "Arctic Dogs",
		Genre:       "Animated",
		MinAge:      1,
		ReleaseDate: time.Date(2019, 11, 12, 0, 0, 0, 0, time.UTC),
	},
}

var AllAgesMovies = []Movie{}

func CreatePreferences(genre string, ageAppropriate int, date time.Time) *Preferences {

	SelectedPreferences := Preferences{
		Genre:          genre,
		AgeAppropriate: ageAppropriate,
		Date:           date,
	}
	return &SelectedPreferences
}

func (g *Preferences) Recommendation(SelectedPreferences *Preferences) (Movie, bool) {
	for _, movie := range movies {

		if movie.Genre == g.Genre && movie.MinAge <= g.AgeAppropriate && movie.ReleaseDate.Before(g.Date) {
			return movie, true
		}
	}
	return Movie{}, false
}

func (g *Preferences) RecommendationForGivenAge(SelectedPreferences *Preferences) Movie {
	//insert SQL request WHERE these conditions are met
	for _, movie := range movies {
		if movie.MinAge <= g.AgeAppropriate && movie.Genre != g.Genre {
			return movie
		}

	}
	return Movie{}
}

func SetMovie(a []Movie) {
	movies = a
}
func (a *MovieService) AddMovie(titleInput string, genreInput string, minAgeInput int) error {

	_, err := a.db.Exec(insertMovieQuery, titleInput, genreInput, minAgeInput)
	if err != nil {
		return err
	}

	return nil
}
func ListMovies() []Movie {
	return movies
}
