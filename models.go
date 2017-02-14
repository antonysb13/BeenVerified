package main

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