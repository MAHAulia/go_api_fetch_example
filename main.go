package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	blacklist = make(map[string]struct{})
	mu        sync.Mutex
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(DefaultHeaders)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Welcome to API Consumer"}`))
	})
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Hello, world!"}`))
	})

	r.Get("/todos", GetTodos)
	r.Get("/todos/{id}", GetDetailTodos)

	fmt.Println("App running on localhost:3000")
	http.ListenAndServe(":3000", r)
}

func DefaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set default Content-Type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func GetDetailTodos(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	response, err := http.Get(fmt.Sprintf(`https://jsonplaceholder.typicode.com/todos/%s`, id))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf(`{"message": "Success", "data": %s}`, responseData)))
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://jsonplaceholder.typicode.com/todos")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf(`{"message": "Success", "data": %s}`, responseData)))
}
