package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	Id string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)

	for _, item := range books {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = strconv.Itoa(rand.Intn(1000000))

	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	params := mux.Vars(r)
	
	var bookData Book
	_ = json.NewDecoder(r.Body).Decode(&bookData)

	for index, item := range books {
		if item.Id == params["id"] {
			
			books = append(books[:index], books[index:+1]...)

			json.NewEncoder(w).Encode(item)
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	params := mux.Vars(r)
	
	var bookData Book
	_ = json.NewDecoder(r.Body).Decode(&bookData)

	for ix, item := range books {
		if item.Id == params["id"] {
			
			item.Author = bookData.Author
			item.Title = bookData.Title
			item.Isbn = bookData.Isbn

			books[ix] = item

			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	books = append(books, Book {Id: "1", Isbn: "jkhkg", Title: "Book1", Author: &Author {
		FirstName: "Paolo", LastName: "Posso",
	}})
	books = append(books, Book {Id: "2", Isbn: "123486", Title: "Book2", Author: &Author {
		FirstName: "Jorge", LastName: "Barbosa",
	}})

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}