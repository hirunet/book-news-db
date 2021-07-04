package main

import (
	"fmt"
  "os"
  "time"
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
	  logMessage(fmt.Sprintf("%v", err))
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
	  logMessage(fmt.Sprintf("%v", err))
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
	  logMessage(fmt.Sprintf("%v", err))
    return
  }

}
func UpdateBook(book BookInfo) {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")

  db, err := sql.Open("postgres", url)
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
    return
  }
  defer db.Close()

  _, err = db.Exec("UPDATE books SET title = $1, author = $2, publisher = $3, pubdate = $4, cover = $5, keywords = $6, ccode = $7, genre = $8, description = $9, contents = $10 WHERE isbn = $11;",
    book.Title,
    book.Author,
    book.Publisher,
    book.Pubdate,
    book.Cover,
    book.Keywords,
    book.Ccode,
    book.Genre,
    book.Description,
    book.Contents,
    book.Isbn)
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
    return
  }

}

func GetIsbnList() []string {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")
  var isbnList []string

  db, err := sql.Open("postgres", url)
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
  }
  defer db.Close()

  rows, err := db.Query("SELECT isbn FROM books;")
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
  }
  defer rows.Close()

  for rows.Next() {
    var isbn string
    err = rows.Scan(&isbn)
    if err != nil {
	    logMessage(fmt.Sprintf("%v", err))
    }
    isbnList = append(isbnList, isbn)
  }

  return isbnList
}


// 更新対象のISBNのリストを取得する
func GetIsbnListToUpdate() []string {
  // TODO sslmode
  var url = os.Getenv("DATABASE_URL")
  var isbnList []string
  var now = time.Now()
  var fromDate = now.AddDate(0, 0, -14)
  var toDate = now.AddDate(0, 1, 0)

  db, err := sql.Open("postgres", url)
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
  }
  defer db.Close()

  rows, err := db.Query("SELECT isbn FROM books WHERE $1 <= pubdate AND pubdate <= $2;", fromDate.Format("20060102"), toDate.Format("20060102"))
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
  }
  defer rows.Close()

  for rows.Next() {
    var isbn string
    err = rows.Scan(&isbn)
    if err != nil {
	    logMessage(fmt.Sprintf("%v", err))
    }
    isbnList = append(isbnList, isbn)
  }

  return isbnList
}

