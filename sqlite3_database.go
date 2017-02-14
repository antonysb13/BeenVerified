package main

import (
	"fmt"

	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

/* Constants */

//File path of the database
const databaseFilePath = "./jrdd.db"

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