package handlers

import (
	"net/http"
)

func ReadEntry(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from ReadEntry"))
}
