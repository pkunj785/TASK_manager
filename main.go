package main

import (
	"net/http"

	"com.serve_volt/auth"
	"com.serve_volt/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/registor", auth.RegisterUser)
	r.Post("/login", auth.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticateUser)
		router.Routers(r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route dosen't exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	http.ListenAndServe(":8080", r)

}
