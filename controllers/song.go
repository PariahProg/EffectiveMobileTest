package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"EffectiveMobileTest/entities"
	"EffectiveMobileTest/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// @Summary Get lyrics of a song
// @Description Get lyrics of a song by id with pagination
// @Tags songs
// @Produce  json
// @Param id path int true "Song id"
// @Param page query int true "Page number"
// @Param versesPerPage query int true "Number of verses per page"
// @Success 200 {object} entities.SongVerses "Successfully fetched lyrics"
// @Failure 400 {string} string "One of parameters is invalid or not provided"
// @Failure 404 {string} string "No song with such id"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id}/lyrics [get]
func GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get song lyrics request received")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		logrus.WithField("id", vars["id"]).Warn("Invalid song id provided")
		http.Error(w, "Invalid song id!", http.StatusBadRequest)
		return
	}

	logrus.WithField("song_id", id).Debug("Fetching song lyrics")

	pageStr := r.URL.Query().Get("page")
	versesPerPageStr := r.URL.Query().Get("versesPerPage")
	if pageStr == "" {
		logrus.Warn("Page parameter not provided")
		http.Error(w, "page parameter not provided!", http.StatusBadRequest)
		return
	}
	if versesPerPageStr == "" {
		logrus.Warn("versesPerPage parameter not provided")
		http.Error(w, "versesPerPage parameter not provided!", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logrus.WithField("page", pageStr).Warn("Invalid page parameter provided")
		http.Error(w, "Incorrect page number!", http.StatusBadRequest)
		return
	}

	versesPerPage, err := strconv.Atoi(versesPerPageStr)
	if err != nil || versesPerPage < 1 {
		logrus.WithField("versesPerPage", versesPerPageStr).Warn("Invalid versesPerPage parameter provided")
		http.Error(w, "Incorrect versesPerPage number!", http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"page":          page,
		"versesPerPage": versesPerPage,
	}).Debug("Parsed query parameters successfully")

	songVerses := entities.SongVerses{Page: page, VersesPerPage: versesPerPage}
	songVerses.Verses, err = models.GetSongLyrics(id, page, versesPerPage)
	if err != nil && err == models.ErrNoSongFound {
		logrus.WithField("id", id).Warn("No song with provided id")
		http.Error(w, "No song with such id!", http.StatusNotFound)
		return
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
			"error":   err,
		}).Error("Error fetching song lyrics")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
			"page":    page,
			"verses":  len(songVerses.Verses),
		}).Info("Fetched song lyrics successfully")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(&songVerses)
		if err != nil {
			logrus.WithField("error", err).Error("Error encoding response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		logrus.WithFields(logrus.Fields{
			"song_id":       id,
			"page":          page,
			"versusPerPage": versesPerPage,
		}).Info("Response successfully sent")
	}
}

// @Summary Add a new song
// @Description Adds a new song to the library. The request body must be in JSON format and include the song's title and group
// @Accept  json
// @Param song body entities.Song true "Song object containing title and group"
// @Success 201 "Song created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 415 {string} string "Unsupported Media Type"
// @Failure 422 {string} string "Incorrect song data provided or has invalid format"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [post]
func AddSong(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Add song request received")
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		logrus.WithField("Content-Type", ct).Warn("Unsupported Content-Type provided")
		http.Error(w, "Unsupported Content-Type!", http.StatusUnsupportedMediaType)
		return
	}

	var song entities.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithField("err", err).Error("Decoding body JSON error")
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}

	if song.Title == "" || song.Group == "" {
		logrus.WithFields(logrus.Fields{
			"title": song.Title,
			"group": song.Group,
		}).Warn("Invalid song data")
		http.Error(w, "Incorrect data provided!\nJSON should contain title and group!", http.StatusUnprocessableEntity)
		return
	}

	api := os.Getenv("API_URL")
	url := fmt.Sprintf("http://%s/info?group=%s&song=%s", api, url.QueryEscape(song.Group), url.QueryEscape(song.Title))
	logrus.WithFields(logrus.Fields{
		"title": song.Title,
		"group": song.Group,
		"url":   url,
	}).Debug("Fetching song data from side API")
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"group": song.Group,
			"title": song.Title,
			"error": err,
		}).Error("Error fetching song data from side API")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"group":  song.Group,
			"title":  song.Title,
			"status": resp.StatusCode,
		}).Error("Side API responsed with error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"group": song.Group,
			"title": song.Title,
			"error": err,
		}).Error("Error decoding response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	releaseDate, err := time.Parse("02.01.2006", song.ReleaseDate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"group": song.Group,
			"title": song.Title,
			"error": err,
		}).Error("Error parsing date from API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	song.ReleaseDate = releaseDate.Format("2006-01-02")

	err = models.AddSong(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"group": song.Group,
			"title": song.Title,
			"error": err,
		}).Error("Error adding song to database")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"group": song.Group,
		"title": song.Title,
	}).Info("Song successfully added")
	w.WriteHeader(http.StatusCreated)
}

// @Summary Update an existing song
// @Description Update the details of a song by its id. The request body must be in JSON format and include all required fields: title, group, releaseDate, lyrics and link.
// @Tags songs
// @Accept json
// @Param id path int true "Song id"
// @Param song body entities.Song true "Song object that needs to be updated"
// @Success 204 "Successfully updated"
// @Failure 400 {string} string "Invalid request body or song id"
// @Failure 404 {string} string "No song found with the provided id"
// @Failure 415 {string} string "Unsupported Content-Type"
// @Failure 422 {string} string "Incorrect body data provided or has invalid format"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [put]
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Updated song request received")
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		logrus.WithField("Content-Type", ct).Warn("Unsupported Content-Type provided")
		http.Error(w, "Unsupported Content-Type!", http.StatusUnsupportedMediaType)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	logrus.WithField("song_id", idStr).Debug("song id parsed")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logrus.WithField("id", vars["id"]).Warn("Invalid song id provided")
		http.Error(w, "Invalid id!", http.StatusBadRequest)
		return
	}

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		logrus.WithField("err", err).Error("Decoding body JSON error")
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}

	if song.Title == "" || song.Group == "" || song.ReleaseDate == "" || song.Lyrics == "" || song.Link == "" {
		logrus.WithFields(logrus.Fields{
			"title":       song.Title,
			"group":       song.Group,
			"releaseDate": song.ReleaseDate,
			"lyrics":      song.Lyrics,
			"link":        song.Link,
		}).Warn("Invalid song data")
		http.Error(w, "Incorrect data provided!\nJSON should contain title, group, release_date, lyrics and link!", http.StatusUnprocessableEntity)
		return
	}

	releaseDate, err := time.Parse("02.01.2006", song.ReleaseDate)
	if err != nil {
		logrus.WithField("releaseDate", song.ReleaseDate).Warn("Invalid release date provided")
		http.Error(w, "Incorrect releaseDate! Please use format DD.MM.YYYY!", http.StatusUnprocessableEntity)
		return
	}
	song.ReleaseDate = releaseDate.Format("2006-01-02")

	logrus.WithFields(logrus.Fields{
		"song_id":     id,
		"title":       song.Title,
		"group":       song.Group,
		"releaseDate": song.ReleaseDate,
		"link":        song.Link,
	}).Debug("Trying update song")

	err = models.UpdateSong(id, &song)
	if err != nil && errors.Is(err, models.ErrNoSongFound) {
		logrus.WithField("song_id", id).Warn("No song with provided id")
		http.Error(w, "No song with such id!", http.StatusNotFound)
		return
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
			"error":   err,
		}).Error("Error updating song")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		logrus.WithField("song_id", id).Info("Song updated successfully")
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Patch song
// @Description Partially updates an existing song by its id. At least one field (title, group, releaseDate, lyrics, or link) must be provided for update.
// @Tags songs
// @Accept json
// @Param id path int true "Song id"
// @Param song body entities.Song true "Fields to update in the song. At least one of: title, group, releaseDate, lyrics, or link."
// @Success 204 "Successfully patched"
// @Failure 400 {string} string "Invalid request body or song id"
// @Failure 415 {string} string "Unsupported Content-Type"
// @Failure 422 {string} string "Incorrect body data provided or has invalid format"
// @Failure 404 {string} string "No song with such id"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [patch]
func PatchSong(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Patching song request received")

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		logrus.WithField("Content-Type", ct).Warn("Unsupported Content-Type provided")
		http.Error(w, "Unsupported Content-Type!", http.StatusUnsupportedMediaType)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	logrus.WithField("song_id", idStr).Debug("song id parsed")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logrus.WithField("id", idStr).Warn("Invalid song id provided")
		http.Error(w, "Invalid id!", http.StatusBadRequest)
		return
	}

	var song entities.Song
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithField("err", err).Error("Decoding body JSON error")
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}

	if song.Title == "" && song.Group == "" && song.ReleaseDate == "" && song.Lyrics == "" && song.Link == "" {
		logrus.Warn("No data for patching provided")
		http.Error(w, "Incorrect data provided!\nJSON should contain at least title, group, releaseDate, lyrics or link!", http.StatusUnprocessableEntity)
		return
	}

	if song.ReleaseDate != "" {
		releaseDate, err := time.Parse("02.01.2006", song.ReleaseDate)
		if err != nil {
			logrus.WithField("releaseDate", song.ReleaseDate).Warn("Invalid releaseDate provided")
			http.Error(w, "Incorrect releaseDate! Please use format DD.MM.YYYY!", http.StatusUnprocessableEntity)
			return
		}
		song.ReleaseDate = releaseDate.Format("2006-01-02")
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     id,
		"title":       song.Title,
		"group":       song.Group,
		"releaseDate": song.ReleaseDate,
		"lyrics":      song.Lyrics,
		"link":        song.Link,
	}).Debug("Trying patching song")

	err = models.PatchSong(id, &song)
	if err != nil && errors.Is(err, models.ErrNoSongFound) {
		logrus.WithField("song_id", id).Warn("No song with provided id")
		http.Error(w, "No song with such id!", http.StatusNotFound)
		return
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
			"error":   err,
		}).Error("Error patching song")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		logrus.WithField("song_id", id).Info("Song successfully patched")
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Delete a song
// @Description Delete a song by id
// @Tags songs
// @Param id path int true "Song id"
// @Success 204 "Successfully deleted"
// @Failure 400 {string} string "Invalid song id"
// @Failure 404 {string} string "No song with such id"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Deleting song request received")

	vars := mux.Vars(r)
	idStr := vars["id"]
	logrus.WithField("song_id", idStr).Debug("song id parsed")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logrus.WithField("id", idStr).Warn("Invalid song id provided")
		http.Error(w, "Invalid id!", http.StatusBadRequest)
		return
	}

	logrus.WithField("song_id", id).Debug("Trying deleting song")

	err = models.DeleteSong(id)
	if err != nil && errors.Is(err, models.ErrNoSongFound) {
		logrus.WithField("song_id", id).Warn("No song with provided id")
		http.Error(w, "No song with such id!", http.StatusNotFound)
		return
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
			"error":   err,
		}).Error("Error deleting song")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		logrus.WithField("song_id", id).Info("Song successfully deleted")
		w.WriteHeader(http.StatusNoContent)
	}
}
