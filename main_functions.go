package main

import (
	"fmt"
	"encoding/json"

	"net/http"
	"goji.io/pat"

	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

//findAllSongs finds all the songs in the database
func findAllSongs(w http.ResponseWriter, r *http.Request){

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get all songs in database
    rows := findAllSongsDB(database)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    printResultAsJSON(w, rows)
}

//findSongByArtist finds all the songs in the database that match with the given artist
func findSongByArtist(w http.ResponseWriter, r *http.Request){

	//Get the parameter value
	artist := pat.Param(r, "artist")

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get the songs in database that match with the given artist
    rows := findSongByArtistDB(database, artist)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    printResultAsJSON(w, rows)
}

//findSongBySong finds all the songs in the database that match with the given song
func findSongBySong(w http.ResponseWriter, r *http.Request){

	//Get the parameter value
	song := pat.Param(r, "song")

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get the songs in database that match with the given song
    rows := findSongBySongDB(database, song)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    printResultAsJSON(w, rows)
}

//findSongByGenre finds all the songs in the database that match with the given genre
func findSongByGenre(w http.ResponseWriter, r *http.Request){

	//Get the parameter value
	genre := pat.Param(r, "genre")

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get the songs in database that match with the given genre
    rows := findSongByGenreDB(database, genre)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    printResultAsJSON(w, rows)
}

//findSongByLength finds all the songs in the database that have a length between a minimum and maximum
func findSongByLength(w http.ResponseWriter, r *http.Request){

	//Get the parameter values
	minLength := pat.Param(r, "minLength")
	maxLength := pat.Param(r, "maxLength")

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get the songs in database that match with the given genre
    rows := findSongByLengthDB(database, minLength, maxLength)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    printResultAsJSON(w, rows)
}

//findAllGenres finds all the genres in the database and gives the number of songs and the total length of all songs by genre
func findAllGenres(w http.ResponseWriter, r *http.Request){

	//Initilize and open the database
	database := initDatabase(databaseFilePath)
	defer database.Close()

	//Get all songs in database
    rows := findAllGenresDB(database)
    defer rows.Close()
 
    //Output the resulted rows as JSON data
    //Encode rows into JSON data
    jsonResponse := genreRowsToJSON(rows)

    //Format the JSON into string and write the result to w  
    fmt.Fprintf(w, string(jsonResponse))
}


//printResultAsJSON outputs the resulted rows as JSON data
func printResultAsJSON(w http.ResponseWriter, rows *sql.Rows){
	//Encode rows into JSON data
    jsonResponse := songRowsToJSON(rows)

    //Format the JSON into string and write the result to w  
    fmt.Fprintf(w, string(jsonResponse))
}

//songRowsToJSON encodes the given rows into JSON data
func songRowsToJSON(rows *sql.Rows) []byte{

	songs := []Song {}

    //Iterate over the rows
    for rows.Next(){

    	song := Song{}

    	songError := rows.Scan(
    		&song.ID, 
    		&song.Artist,
    		&song.Song,
    		&song.Genre,
    		&song.Length)

    	if songError != nil{
    		fmt.Println("Something went wrong trying to get a song.")
        	fmt.Println(songError)
    	}

    	songs = append(songs, song)
    }

    songsListResult := SongsList {
    	Songs: songs,
    }

    //Encode the Go array into JSON data
    jsonResponse,_ := json.Marshal(songsListResult)

    return jsonResponse
}

//genreRowsToJSON encodes the given rows into JSON data
func genreRowsToJSON(rows *sql.Rows) []byte{

	genres := []Genre {}

    //Iterate over the rows
    for rows.Next(){

    	genre := Genre{}

    	genreError := rows.Scan(
    		&genre.Genre, 
    		&genre.NumberOfSongs,
    		&genre.TotalLength)

    	if genreError != nil{
    		fmt.Println("Something went wrong trying to get a genre.")
        	fmt.Println(genreError)
    	}

    	genres = append(genres, genre)
    }

    genresListResult := GenresList {
    	Genres: genres,
    }

    //Encode the Go array into JSON data
    jsonResponse,_ := json.Marshal(genresListResult)

    return jsonResponse
}