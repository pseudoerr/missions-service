package service

import (
	"context"
	"github.com/pseudoerr/mission-service/models"
	"log/slog"
	"sync"
)

type MissionStore interface {
	ListMissions(ctx context.Context) ([]models.Mission, error)
	AddMission(ctx context.Context, m models.Mission) (models.Mission, error)
}

type InMemoryStore struct {
	mu       sync.Mutex
	missions []models.Mission
	nextID   int
}

type MissionService struct {
	Store  MissionStore
	Logger *slog.Logger
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		missions: []models.Mission{
			{ID: 1, Title: "Hello, World!", Points: 100},
			{ID: 2, Title: "FizzBuzz", Points: 200},
		},
		nextID: 3,
	}
}

func (s *InMemoryStore) ListMissions(ctx context.Context) ([]models.Mission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.Mission{}, s.missions...), nil
}

func (s *InMemoryStore) AddMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextID
	s.nextID++
	s.missions = append(s.missions, m)
	return m, nil
}
