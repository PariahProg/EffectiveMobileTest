/*	Файл реализующий подключение к бд.
	Изначально подключение происходит к служебной базе postgres, так как postgresql не разрешает подключаться к серверу без указания конкретной бд.
	После происходит проверка на наличие необходимой бд и, если она отсутствует, происходит создание с последующим заполнением с помощью миграций.
*/

package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	dbUser     string
	dbPassword string
	dbName     string
	dbSsl      string
	Db         *sql.DB
)

func loadEnvVariables() {
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbSsl = os.Getenv("DB_SSL")
}

func isDbExists() (bool, error) {
	var isExists bool
	err := Db.QueryRow("select exists(select * from pg_catalog.pg_database where datname = $1)", dbName).Scan(&isExists)
	return isExists, err
}

func createDb() error {
	_, err := Db.Exec(fmt.Sprintf("create database %s;", dbName))
	return err
}

func applyMigrations() error {
	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil {
		return err
	}
	return nil
}

func OpenDb() error {
	loadEnvVariables()
	connStr := "user=" + string(dbUser) + " password=" + string(dbPassword) + " dbname=postgres" + " sslmode=" + string(dbSsl)
	var err error = nil
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	logrus.Debug("Connection to db server established")

	isExists, err := isDbExists()
	if err != nil {
		return err
	}

	if !isExists {
		logrus.Debug("Db doesn't exist. Trying to create")
		err = createDb()
		if err != nil {
			return err
		}
		connStr = "user=" + string(dbUser) + " password=" + string(dbPassword) + " dbname=" + string(dbName) + " sslmode=" + string(dbSsl)
		Db, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
		logrus.Debug("Db created. Trying to migrate")
		err = applyMigrations()
		if err != nil {
			return err
		}
		logrus.Debug("Migrations applied")
	} else {
		connStr = "user=" + string(dbUser) + " password=" + string(dbPassword) + " dbname=" + string(dbName) + " sslmode=" + string(dbSsl)
		Db, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
	}
	logrus.Debug("Connection to db established")
	return nil
}
