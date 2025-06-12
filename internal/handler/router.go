package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(handler *Handler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/missions", handler.GetMissions).Methods("GET")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.GetMissionByID).Methods("GET")
	r.HandleFunc("/missions", handler.CreateMission).Methods("POST")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.UpdateMission).Methods("PUT")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.DeleteMission).Methods("DELETE")
	r.HandleFunc("/profile", handler.GetProfile).Methods("GET")

	var handlerWithMiddleware http.Handler = r
	handlerWithMiddleware = LoggingMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = RecoverMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = CORSMiddleware(handlerWithMiddleware)

	return handlerWithMiddleware
}
