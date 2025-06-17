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

// GetMissions godoc
// @Summary Получить все задания
// @Description Возвращает список всех заданий
// @Tags missions
// @Produce json
// @Success 200 {array} models.Mission
// @Failure 504 {string} string "Request timeout"
// @Failure 500 {string} string "Failed to list missions"
// @Router /missions [get]
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

// GetMissionByID godoc
// @Summary Получить задание по ID
// @Description Возвращает одно задание по его идентификатору
// @Tags missions
// @Produce json
// @Param id path int true "ID задания"
// @Success 200 {object} models.Mission
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Not found"
// @Router /missions/{id} [get]
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

// CreateMission godoc
// @Summary Создать новое задание
// @Description Добавляет новое задание в систему
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body models.Mission true "Новое задание"
// @Success 201 {object} models.Mission
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create mission"
// @Router /missions [post]
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

// UpdateMission godoc
// @Summary Обновить задание
// @Description Обновляет существующее задание по ID
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID задания"
// @Param mission body models.Mission true "Обновленные данные"
// @Success 200 {object} models.Mission
// @Failure 400 {string} string "Invalid ID or request"
// @Failure 500 {string} string "Failed to update mission"
// @Router /missions/{id} [put]
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

// DeleteMission godoc
// @Summary Удалить задание
// @Description Удаляет задание по ID
// @Tags missions
// @Param id path int true "ID задания"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete mission"
// @Router /missions/{id} [delete]
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

// GetProfile godoc
// @Summary Получить профиль
// @Description Возвращает профиль текущего пользователя/сервиса
// @Tags profile
// @Produce json
// @Success 200 {object} models.Profile
// @Failure 500 {string} string "Failed to get profile"
// @Router /profile [get]
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
