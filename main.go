package main

import (
	"database/sql"
	"glup3/handlers"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	// config file for railway
	viper.SetConfigFile("ENV")
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8052")

	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile(".env")
		log.Println("Reading from .env file")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
	}

	connStr := viper.GetString("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database %s", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Unable to get database driver %s", err)
	}

	mi, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Unable to locate migrations %s", err)
	}

	log.Println("Applying migrations")
	if err := mi.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Migrations are up to date.")
		} else {
			log.Fatalf("Unable to apply migrations %s", err)
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/trecord", handlers.ReadEntry).Methods("GET")

	PORT := viper.GetString("PORT")
	log.Println("Listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
