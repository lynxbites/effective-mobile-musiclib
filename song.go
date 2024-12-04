package musiclib

import "github.com/jackc/pgx/v5/pgtype"

type Song struct {
	Id          string      `json:"id"`
	Group       string      `json:"group"`
	Name        string      `json:"name"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Text        string      `json:"text"`
	Link        string      `json:"link"`
}

type SongPost struct {
	Group       *string      `json:"group"`
	Name        *string      `json:"name"`
	ReleaseDate *pgtype.Date `json:"releaseDate"`
	Text        *string      `json:"text"`
	Link        *string      `json:"link"`
}
