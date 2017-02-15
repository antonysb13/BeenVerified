## BeenVerified Challenge

This project is the solution to the BeenVerified Challenge.

## Getting Started

Follow the next instuctions to get a copy of the project and run it on your local machine.

### Prerequisites

This API was implemented in [Golang 1.7.5](https://golang.org/dl/).

It is necessary to install Glide to get the next dependencies:
	* github.com/mattn/go-sqlite3 - [Go-SQLite3](https://github.com/mattn/go-sqlite3) 
	* goji.io - [Goji](https://github.com/goji/goji) 
	* golang.org/x/net

See how to install [Glide in this repository](https://github.com/Masterminds/glide).

### Installing

To get a copy of the project in your local machine run the next instruction in your terminal or command line:

```
git clone https://github.com/antonysb13/BeenVerified.git $GOPATH/src/github.com/antonysb13/BeenVerified
```

Move to project's directory:

```
cd $GOPATH/src/github.com/antonysb13/BeenVerified
```

To install the dependencies with glide:

```
glide install
```

To build the project and get the executable file:

```
go build
```

### Running the project 

To run the project execute the file called:``` BeenVerified ```

If you are using a Linux distribution move first to the project's directory in your terminal and use:

```
./BeenVerified
```

## API - List of Routes

The routes to access to the API functions are the next:

### Get all the songs

```
http://localhost:8080/songs
```

### Get songs by artist

```
http://localhost:8080/songs/artist/:artist
```

Put the text you want to search instead of ":artist". For example: http://localhost:8080/songs/artist/beatles

### Get songs by song

```
http://localhost:8080/songs/song/:song
```

Put the text you want to search instead of ":song". For example: http://localhost:8080/songs/song/twist

### Get songs by genre

```
http://localhost:8080/songs/genre/:genre
```

Put the text you want to search instead of ":genre". For example: http://localhost:8080/songs/genre/rock

### Get songs by length

```
http://localhost:8080/songs/length/:minLength/:maxLength
```

Put the minimum and maximum length you want to search instead of ":minLength" and ":maxLength" respectively. 
For example, to get the songs between 200 and 245 length: http://localhost:8080/songs/length/200/245

### Get the list of genres, and the number of songs and the total length of all the songs by genre

```
http://localhost:8080/genres
```

## Author

**Antony Sandoval Bonilla** - [My Github Page](https://github.com/antonysb13/)

