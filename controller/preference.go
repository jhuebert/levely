package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"io/ioutil"
	"net/http"
)

func (c *Controller) registerPreferenceRoutes(router *mux.Router) {
	router.HandleFunc("/api/preference", c.GetPreferences).Methods("GET")
	router.HandleFunc("/api/preference", c.UpdatePreferences).Methods("PUT")
}

func (c *Controller) GetPreferences(w http.ResponseWriter, r *http.Request) {
	p := c.s.GetPreferences()
	sendResponse(w, p)
}

func (c *Controller) UpdatePreferences(w http.ResponseWriter, r *http.Request) {

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
		return repository.Preferences{}, err
	}

	p := repository.Preferences{}
	err = json.Unmarshal(body, &p)
	return p, err
}
