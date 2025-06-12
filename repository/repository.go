package repository

import (
	"context"
	"database/sql"

	"github.com/pseudoerr/mission-service/models"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) ListMissions(ctx context.Context) ([]models.Mission, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, title, points FROM missions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var m models.Mission
		if err := rows.Scan(&m.ID, &m.Title, &m.Points); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *PostgresRepository) AddMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	err := r.DB.QueryRowContext(
		ctx,
		"INSERT INTO missions (title, points) VALUES ($1, $2) RETURNING id",
		m.Title, m.Points,
	).Scan(&m.ID)

	return m, err
}
