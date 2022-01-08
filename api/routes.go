package api

import (
	"echo-http/pkg/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(Handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(Handler.JsonAllInfo)
	//r.HandleFunc("/", Handler.JsonAllInfo)
	//r.HandleFunc("/mounts", Handler.JsonMounts).Methods("GET")
	//r.HandleFunc("/env", Handler.JsonEnvs).Methods("GET")
	return r
}
