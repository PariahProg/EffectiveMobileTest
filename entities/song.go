package entities

type Song struct {
	Id          int
	Title       string `json:"title"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}

type SongVerses struct {
	Verses        []string `json:"verses"`
	Page          int      `json:"page"`
	VersesPerPage int      `json:"versesPerPage"`
}
