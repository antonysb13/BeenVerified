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
	ID int
	Name string
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
	//mux.HandleFunc(pat.Get("/songs/genre/:genre"), findSongByGenre)
	/*mux.HandleFunc(pat.Get("/songs/length/:minLength/:maxLength"), findSongByLength)
	
	//Genres Handlers
	mux.HandleFunc(pat.Get("/genres"), findAllGenres)
	*/
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

//printResultAsJSON outputs the resulted rows as JSON data
func printResultAsJSON(w http.ResponseWriter, rows *sql.Rows){
	//Encode rows into JSON data
    jsonResponse := songRowsToJSON(rows)

    fmt.Println(string(jsonResponse))

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

    //Encode the Go object into JSON data
    jsonResponse,_ := json.Marshal(songsListResult)

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

	//Prepare the sql statement
	sqlStmtPrepared, sqlStmtError := database.Prepare(sqlStatement)

	if sqlStmtError != nil {
		fmt.Println("Something went wrong preparing the sql statement: " + sqlStatement)
        fmt.Println(sqlStmtError)
	}
	defer sqlStmtPrepared.Close()

    //Execute the sql statement
    rows, rowsError := sqlStmtPrepared.Query()

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }

    return rows
}

//findSongByArtistDB gets the songs in database that match with the given artist
func findSongByArtistDB(database *sql.DB, artist string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM (SELECT * FROM songs WHERE artist LIKE ?) as S" + 
																		" INNER JOIN genres as G on S.genre = G.ID"

	//Prepare the sql statement
	sqlStmtPrepared, sqlStmtError := database.Prepare(sqlStatement)

	if sqlStmtError != nil {
		fmt.Println("Something went wrong preparing the sql statement: " + sqlStatement)
        fmt.Println(sqlStmtError)
	}
	defer sqlStmtPrepared.Close()

	parameter := "%" + artist + "%"

    //Execute the sql statement
    rows, rowsError := sqlStmtPrepared.Query(parameter)

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }

    return rows
}

//findSongBySongDB gets the songs in database that match with the given song
func findSongBySongDB(database *sql.DB, song string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM (SELECT * FROM songs WHERE song LIKE ?) as S" + 
																		" INNER JOIN genres as G on S.genre = G.ID"

	//Prepare the sql statement
	sqlStmtPrepared, sqlStmtError := database.Prepare(sqlStatement)

	if sqlStmtError != nil {
		fmt.Println("Something went wrong preparing the sql statement: " + sqlStatement)
        fmt.Println(sqlStmtError)
	}
	defer sqlStmtPrepared.Close()

	parameter := "%" + song + "%"

    //Execute the sql statement
    rows, rowsError := sqlStmtPrepared.Query(parameter)

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }

    return rows
}

//findSongByGenreDB gets the songs in database that match with the given genre
func findSongByGenreDB(database *sql.DB, genre string) *sql.Rows{
	sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM songs as S INNER JOIN genres as G on S.genre = G.ID " + 
																					"WHERE G.name LIKE ?"

	//Prepare the sql statement
	sqlStmtPrepared, sqlStmtError := database.Prepare(sqlStatement)

	if sqlStmtError != nil {
		fmt.Println("Something went wrong preparing the sql statement: " + sqlStatement)
        fmt.Println(sqlStmtError)
	}
	defer sqlStmtPrepared.Close()

	parameter := "%" + genre + "%"

    //Execute the sql statement
    rows, rowsError := sqlStmtPrepared.Query(parameter)

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }

    return rows
}

//



