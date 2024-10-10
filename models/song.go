package models

import (
	"EffectiveMobileTest/entities"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrNoSongFound = errors.New("no song found with provided id")

func AddSong(song *entities.Song) error {
	_, err := Db.Exec("INSERT INTO songs (title, group_name, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5)",
		song.Title, song.Group, song.ReleaseDate, song.Lyrics, song.Link)
	if err != nil {
		return fmt.Errorf("error while adding song: %w", err)
	}
	return nil
}

func UpdateSong(id int, song *entities.Song) error {
	result, err := Db.Exec("UPDATE songs SET title = $1, group_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6",
		song.Title, song.Group, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf("error while updating song: %w", err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while checking affecting rows: %w", err)
	}
	if ra == 0 {
		return ErrNoSongFound
	}
	return nil
}

func PatchSong(id int, song *entities.Song) error {
	result, err := Db.Exec("UPDATE songs SET title = COALESCE(NULLIF($1, ''), title), group_name = COALESCE(NULLIF($2, ''), group_name), release_date = COALESCE(NULLIF($3, '')::date, release_date), lyrics = COALESCE(NULLIF($4, ''), lyrics), link = COALESCE(NULLIF($5, ''), link) WHERE id = $6",
		song.Title, song.Group, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf("error while patching song: %w", err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while checking affecting rows: %w", err)
	}
	if ra == 0 {
		return ErrNoSongFound
	}
	return nil
}

func DeleteSong(id int) error {
	result, err := Db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error while deleting song: %w", err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while checking affecting rows: %w", err)
	}
	if ra == 0 {
		return ErrNoSongFound
	}
	return nil
}

func GetSongLyrics(id int, page int, versesPerPage int) ([]string, error) {
	var song entities.Song
	row := Db.QueryRow("SELECT lyrics FROM songs WHERE id = $1", id)
	err := row.Scan(&song.Lyrics)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoSongFound
	} else if err != nil {
		return nil, err
	}

	serverOs := os.Getenv("SERVER_OS")
	var splitter string
	if serverOs == "windows" { // В windows и linux конец строки задается разной последовательностью управляющих символов
		splitter = "\r\n\r\n"
	} else {
		splitter = "\n\n"
	}

	verses := strings.Split(song.Lyrics, splitter)

	start := (page - 1) * versesPerPage
	end := start + versesPerPage
	if start >= len(verses) {
		return []string{}, nil
	}
	if end >= len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}
