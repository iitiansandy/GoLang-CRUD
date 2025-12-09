package routes

import (
	"github.com/coddicted/go-mongo-crud/handlers"
	"github.com/gorilla/mux"
)

func Register() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	return r
}
