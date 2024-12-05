package musiclib

type Song struct {
	Id          string `json:"id"`
	Group       string `json:"group"`
	Name        string `json:"name"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongPaginated struct {
	Id          string   `json:"id"`
	Group       string   `json:"group"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"releaseDate"`
	Text        []string `json:"text"`
	Link        string   `json:"link"`
}

type SongPost struct {
	Group       *string `json:"group"`
	Name        *string `json:"name"`
	ReleaseDate *string `json:"releaseDate"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
}

type SongPatch struct {
	Group       string `json:"group,omitempty"`
	Name        string `json:"name,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}
