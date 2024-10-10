package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"EffectiveMobileTest/controllers"
	_ "EffectiveMobileTest/docs"
	"EffectiveMobileTest/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logrus.SetOutput(os.Stdout)

	env := os.Getenv("APP_ENV")
	if env == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	err := models.OpenDb()
	if err != nil {
		logrus.Fatal("Error connecting to db ", err)
	}
	defer func() {
		if err := models.Db.Close(); err != nil {
			logrus.Error("Error closing db connection", err)
		} else {
			logrus.Info("db connection closed")
		}
	}()

	router := mux.NewRouter()

	router.HandleFunc("/songs", controllers.AddSong).Methods(http.MethodPost) // добавление песни

	router.HandleFunc("/songs/{id:[0-9]+}", controllers.UpdateSong).Methods(http.MethodPut)    // изменение песни
	router.HandleFunc("/songs/{id:[0-9]+}", controllers.PatchSong).Methods(http.MethodPatch)   // частичное изменение песни
	router.HandleFunc("/songs/{id:[0-9]+}", controllers.DeleteSong).Methods(http.MethodDelete) // удаление песни

	router.HandleFunc("/songs/{id:[0-9]+}/lyrics", controllers.GetSongLyrics).Methods(http.MethodGet) // получение текста песни с пагинацией по куплетам

	router.HandleFunc("/library", controllers.GetLibrary).Methods(http.MethodGet) // получение данных библиотеки с фильтрацией по всем полям и пагинацией

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // swagger UI

	port := os.Getenv("SERVER_PORT")
	logrus.Info("Server is started on port ", port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
}
