package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/pseudoerr/mission-service/models"
	"github.com/pseudoerr/mission-service/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Service *service.MissionService
}

func (h *Handler) GetMissions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	missions, err := h.Service.Store.ListMissions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Warn("Mission fetch timed out")
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
		http.Error(w, "Failed to list missions", http.StatusInternalServerError)
		return
	}

	if missions == nil {
		missions = []models.Mission{}
	}

	writeJSON(w, http.StatusOK, missions)
}

func (h *Handler) GetMissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mission, err := h.Service.Store.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, mission)
}

func (h *Handler) CreateMission(w http.ResponseWriter, r *http.Request) {
	var m models.Mission
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	created, err := h.Service.Store.AddMission(r.Context(), m)
	if err != nil {
		http.Error(w, "Failed to create mission", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) UpdateMission(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var m models.Mission
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	m.ID = id

	updated, err := h.Service.Store.UpdateMission(r.Context(), m)
	if err != nil {
		http.Error(w, "Failed to update mission", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteMission(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Service.Store.DeleteMission(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete mission", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.Service.GetProfile(r.Context())
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}
func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}
