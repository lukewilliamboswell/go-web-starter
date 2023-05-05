package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", hanldleGetRoot)

	http.ListenAndServe(":8080", router)
}

func hanldleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
