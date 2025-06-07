package models

type Mission struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Points int    `json:"points"`
}
