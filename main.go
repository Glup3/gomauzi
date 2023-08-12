package main

import (
	"glup3/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/trecord", handlers.ReadEntry).Methods("GET")

	PORT := viper.GetString("PORT")
	log.Println("Listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
