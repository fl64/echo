package api

import (
	"echo-http/internal/app-http/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(Handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/generate").HandlerFunc(Handler.Generate)
	r.PathPrefix("/").HandlerFunc(Handler.JsonAllInfo)
	//r.HandleFunc("/", Handler.JsonAllInfo)
	//r.HandleFunc("/mounts", Handler.JsonMounts).Methods("GET")
	//r.HandleFunc("/env", Handler.JsonEnvs).Methods("GET")
	return r
}
