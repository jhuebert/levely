package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"io/ioutil"
	"net/http"
)

func (c *Controller) registerPositionRoutes(router *mux.Router) {
	router.HandleFunc("/api/position", c.FindAllPositions).Methods("GET")
	router.HandleFunc("/api/position", c.CreatePosition).Methods("POST")
	router.HandleFunc("/api/position/{id}", c.FindPosition).Methods("GET")
	router.HandleFunc("/api/position/{id}", c.UpdatePosition).Methods("PUT")
	router.HandleFunc("/api/position/{id}", c.DeletePosition).Methods("DELETE")
}

func (c *Controller) FindAllPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := c.s.FindAllPositions()
	if err != nil {
		notFoundError(w, err)
		return
	}
	sendResponse(w, positions)
}

func (c *Controller) FindPosition(w http.ResponseWriter, r *http.Request) {
	id := getPathInt(r, "id")
	position, err := c.s.FindPosition(id)
	if err != nil {
		notFoundError(w, err)
		return
	}
	sendResponse(w, position)
}

func (c *Controller) CreatePosition(w http.ResponseWriter, r *http.Request) {

	p, err := readPosition(r)
	if err != nil {
		internalServerError(w, err)
		return
	}

	o, err := c.s.CreatePosition(p)
	if err != nil {
		internalServerError(w, err)
		return
	}

	sendResponse(w, o)
}

func (c *Controller) UpdatePosition(w http.ResponseWriter, r *http.Request) {

	updated, err := readPosition(r)
	if err != nil {
		internalServerError(w, err)
		return
	}

	id := getPathInt(r, "id")

	p, err := c.s.UpdatePosition(id, updated)
	if err != nil {
		notFoundError(w, err) //TODO better status
		return
	}

	sendResponse(w, p)
}

func (c *Controller) DeletePosition(w http.ResponseWriter, r *http.Request) {
	id := getPathInt(r, "id")
	err := c.s.DeletePosition(id)
	if err != nil {
		notFoundError(w, err) //TODO better status
	}
	return
}

func readPosition(r *http.Request) (repository.Position, error) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return repository.Position{}, err
	}

	p := repository.Position{}
	err = json.Unmarshal(body, &p)
	return p, err
}
