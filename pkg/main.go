package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v1/books"
)

type Book struct {
	Id, Name, Isbn string
}

type library struct {
	dbName, dbPass, dbHost string
}

func main() {

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "Pass@123"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	l := library{
		dbName: dbName,
		dbPass: dbPass,
		dbHost: dbHost,
	}

	r := mux.NewRouter()

	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBook).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (l library) postBook(w http.ResponseWriter, r *http.Request) {
	//Read the request into an instance of book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)

	//Open connection
	db := l.openConnection()

	//Write the data
	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("executing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commiting the transaction %s\n", err.Error())
	}

	//Close connection
	l.closeConnection(db)
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	//Open connection
	db := l.openConnection()

	//Read all books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the books table %s\n", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("while scanning the rows %s\n", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}
	json.NewEncoder(w).Encode(books)

	//Close connection
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("opening the connection to database %s\n", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
