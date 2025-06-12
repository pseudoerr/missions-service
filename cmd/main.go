package main

import (
	"database/sql"
	"encoding/json"
	"github.com/pseudoerr/mission-service/repository"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pseudoerr/mission-service/config"
	"github.com/pseudoerr/mission-service/models"
	"github.com/pseudoerr/mission-service/service"
)

type Handler struct {
	Service *service.MissionService
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	config.LoadEnv()
	db, err := sql.Open("postgres", config.GetDatabaseURL())
	if err != nil {
		logger.Warn("failed to connect to db", "postgres", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Warn("failed to ping db", "ping", err)
	}
	logger.Info("connected to psql!")

	repo := repository.NewPostgresRepository(db)
	svc := &service.MissionService{
		Store:  repo,
		Logger: logger,
	}

	handler := &Handler{Service: svc}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/missions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetMissions(w, r)
		case http.MethodPost:
			handler.AddMission(w, r)
		default:
			logger.Warn("Unsupported method", "method", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	logger.Info("Missions API running", "port", 8080)
	http.ListenAndServe(":8080", nil)
}

func (h *Handler) GetMissions(w http.ResponseWriter, r *http.Request) {
	h.Service.Logger.Info("Received request", "path", r.URL.Path, "method", r.Method, "endpoint", "/missions")

	missions, err := h.Service.Store.ListMissions(r.Context())
	if err != nil {
		h.Service.Logger.Error("Error listing missions", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	h.Service.Logger.Info("Missions listed", "count", len(missions))
	_ = json.NewEncoder(w).Encode(missions)
}

func (h *Handler) AddMission(w http.ResponseWriter, r *http.Request) {
	h.Service.Logger.Info("Received request", "path", r.URL.Path, "method", r.Method, "endpoint", "/missions")
	if r.Method != http.MethodPost {
		h.Service.Logger.Warn("Unsupported method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var m models.Mission

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.Service.Logger.Error("Error parsing JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newM, err := h.Service.Store.AddMission(r.Context(), m)
	if err != nil {
		h.Service.Logger.Error("Error adding mission", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	h.Service.Logger.Info("Added mission", "id", newM.ID, "title", newM.Title)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newM)
}
