package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *Handler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/missions", handler.GetMissions).Methods("GET")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.GetMissionByID).Methods("GET")
	r.HandleFunc("/missions", handler.CreateMission).Methods("POST")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.UpdateMission).Methods("PUT")
	r.HandleFunc("/missions/{id:[0-9]+}", handler.DeleteMission).Methods("DELETE")
	r.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	rl := NewRateLimiter(10, time.Minute)

	var handlerWithMiddleware http.Handler = r
	handlerWithMiddleware = rl.MiddleWare(handlerWithMiddleware)
	handlerWithMiddleware = LoggingMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = RecoverMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = CORSMiddleware(handlerWithMiddleware)

	return handlerWithMiddleware
}
