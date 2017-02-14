package main

import (
    "fmt"

    "net/http"

	"goji.io"
	"goji.io/pat"
)

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



