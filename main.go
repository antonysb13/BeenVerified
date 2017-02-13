package main

import (
    "fmt"
    "encoding/json"

    "net/http"

	"goji.io"
	"goji.io/pat"

	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

/* Constants */

//File path of the database
const databaseFilePath = "./jrdd.db"

/* Models */

//Song 
type Song struct{
	ID int 
	Artist string
	Song string
	Genre string
	Length int
}

//Array of Songs
type SongsList struct{
	Songs []Song
}

//Genre
type Genre struct{
	Genre string
	NumberOfSongs int
	TotalLength int
}

//Array of Genres
type GenresList struct{
	Genres []Genre
}

/* Main Function */

func main() {

	fmt.Println("Server starts ...")

	//Handlers
	mux := goji.NewMux()

	//Songs Handlers
	mux.HandleFunc(pat.Get("/songs"), findAllSongs)
	mux.HandleFunc(pat.Get("/songs/artist/:artist"), findSongByArtist)
	mux.HandleFunc(pat.Get("/songs/song/:song"), findSongBySong)
	mux.HandleFunc(pat.Get("/songs/genre/:genre"), findSongByGenre)
	mux.HandleFunc(pat.Get("/songs/length/:minLength/:maxLength"), findSongByLength)
	
	//Genres Handlers
	mux.HandleFunc(pat.Get("/genres"), findAllGenres)
	
	//Host and port of the server
	http.ListenAndServe("localhost:8080", mux)
}

/* Auxiliary Functions */

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

/* Database Functions */

//initDatabase initializes and opens the database located in the given filePath
func initDatabase(filePath string) *sql.DB{
	database, databaseError := sql.Open("sqlite3", filePath)

	if databaseError != nil {
        fmt.Println("Something went wrong openning the database: " + filePath)
        fmt.Println(databaseError)
    }

    return database 
}

//findAllSongsDB gets all songs in database by executing a sql statement
func findAllSongsDB(database *sql.DB) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM songs as S INNER JOIN genres as G on S.genre = G.ID"

	//Execute the query over the database
	rows := executeQuery(database, sqlStatement)

    return rows
}

//findSongByArtistDB gets the songs in database that match with the given artist
func findSongByArtistDB(database *sql.DB, artist string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM (SELECT * FROM songs WHERE artist LIKE ?) as S" + 
																		" INNER JOIN genres as G on S.genre = G.ID"

	parameter := "%" + artist + "%"

	//Execute the query over the database
	rows := executeQuery(database, sqlStatement, parameter)

    return rows
}

//findSongBySongDB gets the songs in database that match with the given song
func findSongBySongDB(database *sql.DB, song string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM (SELECT * FROM songs WHERE song LIKE ?) as S" + 
																		" INNER JOIN genres as G on S.genre = G.ID"

	parameter := "%" + song + "%"

	//Execute the query over the database
    rows := executeQuery(database, sqlStatement, parameter)

    return rows
}

//findSongByGenreDB gets the songs in database that match with the given genre
func findSongByGenreDB(database *sql.DB, genre string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM songs as S INNER JOIN genres as G on S.genre = G.ID " + 
																					"WHERE G.name LIKE ?"
	
	parameter := "%" + genre + "%"

	//Execute the query over the database
	rows := executeQuery(database, sqlStatement, parameter)																				

    return rows
}
 
//findSongByLengthDB gets the songs in database that have a length between a minimum and maximum
func findSongByLengthDB(database *sql.DB, minLength string, maxLength string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM (SELECT * FROM songs WHERE length BETWEEN ? AND ?) as S" + 
																		" INNER JOIN genres as G on S.genre = G.ID"

	//Execute the query over the database
	rows := executeQuery(database, sqlStatement, minLength, maxLength)																				

    return rows
}

//findAllGenresDB gets all genres in database and gives the number of songs and the total length of all songs by genre
func findAllGenresDB(database *sql.DB) *sql.Rows{
	sqlStatement := "SELECT G.name as Genre, COUNT(S.ID) as NumberOfSongs, IFNULL(SUM(S.length), 0 ) as TotalLength FROM genres as G " + 
																		" LEFT OUTER JOIN songs as S on G.ID = S.genre GROUP BY G.name"

	//Execute the query over the database
	rows := executeQuery(database, sqlStatement)

    return rows
}

//executeQuery executes a query over the database with the given parameters 
func executeQuery (database *sql.DB, sqlStatement string, params ...string) *sql.Rows{

	//Prepare the sql statement
	sqlStmtPrepared, sqlStmtError := database.Prepare(sqlStatement)

	if sqlStmtError != nil {
		fmt.Println("Something went wrong preparing the sql statement: " + sqlStatement)
        fmt.Println(sqlStmtError)
	}
	defer sqlStmtPrepared.Close()

	//Execute the sql statement
	var rows *sql.Rows
	var rowsError error
	if len(params) == 1 {
        rows, rowsError = sqlStmtPrepared.Query(params[0])
    }else if len(params) == 2 {
        rows, rowsError = sqlStmtPrepared.Query(params[0], params[1])
    }else{
    	rows, rowsError = sqlStmtPrepared.Query()
    }
    
    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }
    
    return rows
}



