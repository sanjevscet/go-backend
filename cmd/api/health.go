package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK Air !!!"))
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// data := map[string]bool{"ok": true}
	// json.NewEncoder(w).Encode(data)
}
