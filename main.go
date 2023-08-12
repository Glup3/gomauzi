package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	// config file for railway
	viper.SetConfigFile("ENV")

	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile(".env")
		fmt.Println("Reading from .env file")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8052")

	PORT := viper.GetString("PORT")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello GoMauzi"))
	})

	fmt.Println("Server running on port", PORT)

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		panic(err)
	}
}
