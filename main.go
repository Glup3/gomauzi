package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("ENV")
	viper.SetDefault("PORT", "8052")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file")
		panic(err)
	}

	PORT := viper.GetString("PORT")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello GoMauzi"))
	})

	http.ListenAndServe(":"+PORT, nil)
}
