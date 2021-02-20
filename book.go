package main

import (
	"fmt"
  "os"
  _ "github.com/lib/pq"
	"database/sql"
)

type BookInfo struct {
  Isbn        string
  Title       string
  Author      string
  Publisher   string
  Pubdate     string
  Cover       string
  Keywords    string
  Ccode       string
  Genre       string
  Description string
  Contents    string
}

// DBに接続する
func TestDb() int {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")

  db, err := sql.Open("postgres", url)
  if err != nil {
	  fmt.Println(err)
    return 1
  }
  defer db.Close()

  return 0
}

func InsertBook(book BookInfo) {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")

  db, err := sql.Open("postgres", url)
  if err != nil {
	  fmt.Println(err)
    return
  }
  defer db.Close()

  _, err = db.Exec("INSERT INTO books (isbn, title, author, publisher, pubdate, cover, keywords, ccode, genre, description, contents) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
    book.Isbn, 
    book.Title, 
    book.Author, 
    book.Publisher, 
    book.Pubdate, 
    book.Cover, 
    book.Keywords,
    book.Ccode,
    book.Genre,
    book.Description,
    book.Contents)
  if err != nil {
	  fmt.Println(err)
    return
  }

}

func GetIsbnList() []string {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")

  var isbnList []string

  db, err := sql.Open("postgres", url)
  if err != nil {
	  fmt.Println(err)
  }
  defer db.Close()

  rows, err := db.Query("SELECT isbn FROM books;")
  if err != nil {
	  fmt.Println(err)
  }
  defer rows.Close()

  for rows.Next() {
    var isbn string
    err = rows.Scan(&isbn)
    if err != nil {
	    fmt.Println(err)
    }
    isbnList = append(isbnList, isbn)
  }

  return isbnList
}


