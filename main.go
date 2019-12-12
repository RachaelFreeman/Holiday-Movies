package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RachaelFreeman/Holiday-Movies/db"
	"github.com/RachaelFreeman/Holiday-Movies/movie"
	"github.com/RachaelFreeman/Holiday-Movies/storage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manifoldco/promptui"
)

const (
	addMovieCmd  = "Add Movie"
	findMovieCmd = "Find Movie"
)

var movieService *movie.MovieService

func main() {

	db, err := db.ConnectDatabase("movies_db.config")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	movieService = movie.NewService(db)

	err = storage.Load()
	if err != nil {
		fmt.Printf("Failed to Load Movies %v\n", err)
		return
	}
	for {
		fmt.Println("Please select what you would like to do:")

		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				addMovieCmd,
				findMovieCmd,
			},
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case addMovieCmd:
			err := addMovieDetails()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			err = storage.Save()
			if err != nil {
				fmt.Printf("Failed to save movies %v\n", err)
				return
			}

		case findMovieCmd:

			err, genre := selectGenrePrompt()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			ageAppropriate, err := ageAppropriatePrompt()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			date, err := dateOfVist()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			moviePreferences := movie.CreatePreferences(genre, ageAppropriate, date)

			a, b := moviePreferences.Recommendation(moviePreferences)

			printReleaseDate := a.ReleaseDate

			formattedDate := printReleaseDate.Format("January 2")

			movieMinAge := a.MinAge

			rating := ratingConversion(movieMinAge)

			if b == true {
				fmt.Println("You should see", a.Title, "The movie is rated", rating, "and the release date is ", formattedDate, "2019")

			}

			if rating == "PG" {
				fmt.Println("Parental Guidance Recommended")
			}

			releaseDate := a.ReleaseDate
			threeDaysAfter := releaseDate.AddDate(0, 0, 3)

			if date.Before(threeDaysAfter) && date.After(a.ReleaseDate) {
				fmt.Println("This movie is a new release. You may need to buy tickets ahead of time.")
			}

			if b == false {
				fmt.Println("Sorry, there aren't any movies that fit your preferences!")
			}
			moreOptions, err := additionalMovies()
			if moreOptions == "Yes" {
				newMovie := moviePreferences.RecommendationForGivenAge(moviePreferences)
				age := newMovie.MinAge
				newRating := ratingConversion(age)
				newFormattedDate := printReleaseDate.Format("January 2")

				fmt.Println("You may be interested in seeing", newMovie.Title, "The movie is rated", newRating, "and the release date is ", newFormattedDate, "2019")
				if newRating == "PG" {
					fmt.Println("Parental Guidance Recommended")
				}
			}
			time.Sleep(5000 * time.Millisecond)
		}
	}
}

func selectGenrePrompt() (error, string) {

	genrePrompt := promptui.Select{
		Label: "Which genre do you prefer?",
		Items: movie.SelectGenre,
	}

	_, genre, err := genrePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err, ""
	}

	return nil, genre

}

func ageAppropriatePrompt() (int, error) {

	ageAppropriatePrompt := promptui.Prompt{
		Label: "How old is the youngest person in your group?",
	}

	number, err := ageAppropriatePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 0, err
	}

	minAge, err := strconv.Atoi(number)

	return minAge, nil

}

func dateOfVist() (time.Time, error) {
	monthNumber := []string{
		"November",
		"December",
	}

	whenPrompt := promptui.Select{
		Label: "What Month would you like to see a movie?",
		Items: monthNumber,
	}

	_, monthStr, err := whenPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return time.Time{}, err
	}

	var monthInt int
	if monthStr == "November" {
		monthInt = 11
	}

	if monthStr == "December" {
		monthInt = 12
	}

	datePrompt := promptui.Prompt{
		Label: "what day of the month would you like to see the movie? (dd)",
	}

	number, err := datePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return time.Time{}, err
	}

	day, err := strconv.Atoi(number)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return time.Time{}, err
	}

	desiredDate := time.Date(2019, time.Month(monthInt), day, 0, 0, 0, 0, time.UTC)

	return desiredDate, nil

}

func additionalMovies() (string, error) {

	yesNo := []string{
		"Yes",
		"No",
	}
	newPrompt := promptui.Select{
		Label: "Would you like to see additional movie recommendations?",
		Items: yesNo,
	}

	_, result, err := newPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	newSuggestions := result

	return newSuggestions, nil

}
func ratingConversion(minAge int) string {

	if minAge == 0 {
		return "G"
	}

	if minAge == 1 {
		return "PG"
	}
	if minAge == 13 {
		return "PG-13"
	}

	if minAge == 18 {
		return "R"
	}

	return ""
}

func addMovieDetails() error {

	labels := []string{
		"What is the title of the movie?",
		"What is the genre of the movie?",
		"How old do you need to be to watch the movie?",
	}

	movieDetails := []string{}

	for _, labels := range labels {

		titlePrompt := promptui.Prompt{
			Label: labels,
		}
		result, err := titlePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}
		movieDetails = append(movieDetails, result)
	}

	titleInput := movieDetails[0]
	genreInput := movieDetails[1]
	minAgeInput := movieDetails[2]

	minAgeInt, err := strconv.Atoi(minAgeInput)
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	err = movieService.AddMovie(titleInput, genreInput, minAgeInt)

	fmt.Println("Successfuly added", titleInput, "to archive.")
	return nil

}
