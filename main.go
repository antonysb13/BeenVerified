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
	/*mux.HandleFunc(pat.Get("/songs/artist/:artist"), findSongByArtist)
	mux.HandleFunc(pat.Get("/songs/song/:song"), findSongBySong)
	mux.HandleFunc(pat.Get("/songs/genre/:genre"), findSongByGenre)
	mux.HandleFunc(pat.Get("/songs/length/:minLength/:maxLength"), findSongByLength)
	
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
 
    songs := []Song {}

    //Iterate over the returned rows
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

    fmt.Println(string(jsonResponse))

    //Format the JSON into string and write the result to w  
    fmt.Fprintf(w, string(jsonResponse))
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

    //Execute the sql statement
    rows, rowsError := database.Query(sqlStatement)

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement: " + sqlStatement)
        fmt.Println(rowsError)
    }

    return rows
}



