package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Book struct {
	isbn    string
	title   string
	author  string
	price   float32
	created time.Time
	updated time.Time
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://developer:p2ssW0rd@localhost/bookstore")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer rows.Close()

	books := make([]*Book, 0)
	for rows.Next() {
		book := new(Book)
		err := rows.Scan(&book.isbn, &book.title, &book.author, &book.price, &book.created, &book.updated)
		if err = rows.Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, book := range books {
		fmt.Fprintf(w, "%s, %s, %s, $%.2f, %s, %s\n", book.isbn, book.title, book.author, book.price, book.created, book.updated)
	}

}

func booksShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	row := db.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)

	book := new(Book)
	err := row.Scan(&book.isbn, &book.title, &book.author, &book.price, &book.created, &book.updated)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s, %s, %s, $%.2f, %s, %s\n", book.isbn, book.title, book.author, book.price, book.created, book.updated)

}

func booksCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	title := r.FormValue("title")
	author := r.FormValue("author")
	if isbn == "" || title == "" || author == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	created := time.Now()
	updated := time.Now()

	result, err := db.Exec("INSERT INTO books VALUES($1, $2, $3, $4, $5, $6)", isbn, title, author, price, created, updated)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Book %s created successfully (%d row affected)\n", isbn, rowsAffected)

}

func booksUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	title := r.FormValue("title")
	author := r.FormValue("author")
	if isbn == "" || title == "" || author == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	created := time.Now()
	updated := time.Now() //db.QueryRow("SELECT updated FROM books WHERE isbn = $1", isbn).Scan()

	result, err := db.Exec("UPDATE books SET isbn = $1, title = $2, author = $3, price = $4, created = $5, updated = $6 WHERE isbn = $1", isbn, title, author, price, created, updated)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Book %s updated successfully (%d row affected)\n", isbn, rowsAffected)

}

func booksDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	result, err := db.Exec("DELETE FROM books WHERE isbn = $1", isbn)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Book %s deleted successfully (%d row affected)\n", isbn, rowsAffected)

}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/books", booksIndex).Methods("GET")
	api.HandleFunc("/books/show", booksShow).Methods("GET")
	api.HandleFunc("/books/create", booksCreate).Methods("POST")
	api.HandleFunc("/books/update", booksUpdate).Methods("PUT")
	api.HandleFunc("/books/delete", booksDelete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
