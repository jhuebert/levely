package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (c *Controller) registerPreferenceRoutes(router *mux.Router) {
	router.HandleFunc("/html/preference", c.showPreferences).Methods("GET")
	router.HandleFunc("/api/preference", c.updatePreferences).Methods("PUT")
	router.HandleFunc("/api/preference/export", c.exportPreferences).Methods("GET")
}

func (c *Controller) showPreferences(w http.ResponseWriter, r *http.Request) {

	p := c.s.GetPreferences()

	tData := struct {
		Preferences repository.Preferences
	}{
		p,
	}
	err := c.t.ExecuteTemplate(w, "preferences", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) exportPreferences(w http.ResponseWriter, r *http.Request) {
	err := c.s.ExportPreferences(w)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

func (c *Controller) updatePreferences(w http.ResponseWriter, r *http.Request) {

	p, err := readPreferences(r)
	if err != nil {
		internalServerError(w, err)
		return
	}

	p, err = c.s.UpdatePreferences(p)
	if err != nil {
		internalServerError(w, err)
		return
	}

	sendResponse(w, p)
}

func readPreferences(r *http.Request) (repository.Preferences, error) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return repository.Preferences{}, fmt.Errorf("could not read preferences from request body: %w", err)
	}

	p := repository.Preferences{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		err = fmt.Errorf("could not unmarshal preferences from request body: %w", err)
	}

	return p, err
}
