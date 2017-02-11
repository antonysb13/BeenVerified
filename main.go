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

/* Model */

type Song struct{
	ID int 
	Artist string
	Song string
	Genre string
	Length int
}

type SongsList struct{
	Songs []Song
}

type Genre struct{
	ID int
	Name string
}

type GenresList struct{
	Genres []Genre
}

/* Main Function */

func main() {

	fmt.Println("Server starts ...")

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/songs"), findAllSongs)

	http.ListenAndServe("localhost:8080", mux)

}

/* Database Functions */



/* Auxiliary Functions */

func findAllSongs(w http.ResponseWriter, r *http.Request){

	database, databaseError := sql.Open("sqlite3", "./jrdd.db")

	if databaseError != nil {
        fmt.Println("Something went wrong openning the database.")
        fmt.Println(databaseError)
    }
    defer database.Close()

    sqlStatement := "SELECT S.ID, S.artist, S.song, G.name, S.length FROM songs as S INNER JOIN genres as G on S.genre = G.ID"

    rows, rowsError := database.Query (sqlStatement)

    if rowsError != nil {
    	fmt.Println("Something went wrong executing the sql statement.")
        fmt.Println(rowsError)
    }
    defer rows.Close()
 
    songs := []Song {}

    for rows.Next(){
    	var id int
    	var artist string
    	var song string
    	var genre string
    	var length int
    	songError := rows.Scan(
    		&id, 
    		&artist,
    		&song,
    		&genre,
    		&length)

    	if songError != nil{
    		fmt.Println("Something went wrong trying to get a song.")
        	fmt.Println(songError)
    	}

    	songAuxiliar := Song{
    		ID: id,
    		Artist: artist,
    		Song: song,
    		Genre: genre,
    		Length: length}

    	songs = append(songs, songAuxiliar)
    }

    songsListResult := SongsList {
    	Songs: songs,
    }

    jsonResponse,_ := json.Marshal(songsListResult)

    fmt.Println(string(jsonResponse))

    fmt.Fprintf(w, string(jsonResponse))
}

