package models

import (
	"EffectiveMobileTest/entities"
	"fmt"
	"time"
)

func formatSongReleaseDate(song *entities.Song) error {
	releaseDate, err := time.Parse(time.RFC3339, song.ReleaseDate)
	if err != nil {
		return err
	}

	song.ReleaseDate = releaseDate.Format("02.01.2006")
	return nil
}

func GetLibrary(title, group, releaseDate, lyrics, link string, limit, offset int) ([]entities.Song, error) {
	query := `SELECT id, title, group_name, release_date, lyrics, link FROM songs WHERE 1=1`
	args := []interface{}{} // переменная хранит параметры фильтрации и пагинации. Тип переменной []interface{}, так как аргументы имеют типы string и int

	if title != "" {
		query += " AND title ILIKE $" + fmt.Sprint(len(args)+1)
		args = append(args, title+"%")
	}
	if group != "" {
		query += " AND group_name ILIKE $" + fmt.Sprint(len(args)+1)
		args = append(args, group+"%")
	}
	if releaseDate != "" {
		query += " AND release_date = $" + fmt.Sprint(len(args)+1)
		args = append(args, releaseDate)
	}
	if lyrics != "" {
		query += " AND lyrics ILIKE $" + fmt.Sprint(len(args)+1)
		args = append(args, "%"+lyrics+"%")
	}
	if link != "" {
		query += " AND link ILIKE $" + fmt.Sprint(len(args)+1)
		args = append(args, "%"+link+"%")
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	library := []entities.Song{}
	for rows.Next() {
		var song entities.Song
		err := rows.Scan(&song.Id, &song.Title, &song.Group, &song.ReleaseDate, &song.Lyrics, &song.Link)
		if err != nil {
			return nil, err
		}
		formatSongReleaseDate(&song)
		library = append(library, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return library, nil
}
