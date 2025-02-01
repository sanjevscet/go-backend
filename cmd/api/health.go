package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(app.config)

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": API_VERSION,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
