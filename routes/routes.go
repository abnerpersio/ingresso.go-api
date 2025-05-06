package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"ingresso.go/handlers"
)

func Register() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/health", handlers.GetHealth).Methods(http.MethodGet)

	return r
}
