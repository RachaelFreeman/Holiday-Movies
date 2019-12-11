package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/RachaelFreeman/Holiday-Movies/movie"
)

const filename = "movie.json"

func Load() error {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var savedMovies []movie.Movie
	err = json.Unmarshal(fileContents, &savedMovies)
	if err != nil {
		return err
	}

	movie.SetMovie(savedMovies)

	return nil
}

func Save() error {
	moviesList := movie.ListMovies()

	MoviesListBytes, err := json.MarshalIndent(moviesList, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, MoviesListBytes, 0775)
	if err != nil {
		return err
	}

	return nil
}
