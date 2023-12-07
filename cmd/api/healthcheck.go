package main

import (
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountred an error and couldn't process your request", http.StatusInternalServerError)
	}
}
