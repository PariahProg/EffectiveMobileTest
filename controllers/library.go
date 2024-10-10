package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"EffectiveMobileTest/models"

	"github.com/sirupsen/logrus"
)

// @Summary Get songs library
// @Description Retrieve songs from the library with optional filters and pagination
// @Tags library
// @Accept  json
// @Produce  json
// @Param title query string false "Filter by song title"
// @Param group query string false "Filter by group name"
// @Param releaseDate query string false "Filter by release date (format DD.MM.YYYY)"
// @Param lyrics query string false "Filter by lyrics"
// @Param link query string false "Filter by link to clip"
// @Param page query int true "Page number"
// @Param songsPerPage query int true "Number of songs per page"
// @Success 200 {array} entities.Song "Successfully fetched songs library"
// @Failure 400 {string} string "One of query parameters is invalid"
// @Failure 500 {string} string "Internal server error"
// @Router /library [get]
func GetLibrary(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get library request received")
	title := r.URL.Query().Get("title")
	group := r.URL.Query().Get("group")
	releaseDateStr := r.URL.Query().Get("releaseDate")
	lyrics := r.URL.Query().Get("lyrics")
	link := r.URL.Query().Get("link")
	pageStr := r.URL.Query().Get("page")
	songsPerPageStr := r.URL.Query().Get("songsPerPage")

	if pageStr == "" {
		logrus.Warn("Page parameter not provided")
		http.Error(w, "page parameter not provided!", http.StatusBadRequest)
		return
	}
	if songsPerPageStr == "" {
		logrus.Warn("songsPerPage parameter not provided")
		http.Error(w, "songsPerPage parameter not provided!", http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"title":           title,
		"group":           group,
		"releaseDateStr":  releaseDateStr,
		"lyrics":          lyrics,
		"link":            link,
		"pageStr":         pageStr,
		"songsPerPageStr": songsPerPageStr,
	}).Debug("Request to fetch songs library")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		logrus.WithField("pageStr", pageStr).Warn("Invalid page parameter provided")
		http.Error(w, "Invalid page parameter provided!", http.StatusBadRequest)
		return
	}

	songsPerPage, err := strconv.Atoi(songsPerPageStr)
	if err != nil || songsPerPage < 1 {
		logrus.WithField("songsPerPageStr", songsPerPageStr).Warn("Invalid songsPerPage parameter provided")
		http.Error(w, "Invalid songsPerPage parameter provided!", http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"page":         page,
		"songsPerPage": songsPerPage,
	}).Debug("Parsed page and songsPerPage")

	var releaseDate time.Time
	releaseDateFormatted := ""
	if releaseDateStr != "" {
		releaseDate, err = time.Parse("02.01.2006", releaseDateStr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"releaseDateStr": releaseDateStr,
				"err":            err,
			}).Warn("Ivalid releaseDate format")
			http.Error(w, "Invalid releaseDate format! Please use DD.MM.YYYY.", http.StatusBadRequest)
			return
		}
		releaseDateFormatted = releaseDate.Format("2006-01-02")
	}

	limit := songsPerPage
	offset := (page - 1) * songsPerPage

	logrus.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Debug("Calculated limit and offset")

	library, err := models.GetLibrary(title, group, releaseDateFormatted, lyrics, link, limit, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"title":        title,
			"group":        group,
			"releaseDate":  releaseDate,
			"lyrics":       lyrics,
			"link":         link,
			"page":         page,
			"songsPerPage": songsPerPage,
			"err":          err,
		}).Error("Error fetching library data")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"title":        title,
		"group":        group,
		"releaseDate":  releaseDate,
		"lyrics":       lyrics,
		"link":         link,
		"page":         page,
		"songsPerPage": songsPerPage,
	}).Info("Successfully fetched library data")

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&library); err != nil {
		logrus.WithField("error", err).Error("Error encoding response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	logrus.WithFields(logrus.Fields{
		"title":        title,
		"group":        group,
		"releaseDate":  releaseDate,
		"lyrics":       lyrics,
		"link":         link,
		"page":         page,
		"songsPerPage": songsPerPage,
	}).Info("Response successfully returned")
}
