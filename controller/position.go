package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (c *Controller) registerPositionRoutes(router *mux.Router) {
	router.HandleFunc("/position", c.showPositionList).Methods("GET")
	router.HandleFunc("/position/new", c.showNewPosition).Methods("GET")
	router.HandleFunc("/position/{id}", c.showExistingPosition).Methods("GET")
	router.HandleFunc("/api/position", c.createPosition).Methods("POST")
	router.HandleFunc("/api/position/{id}", c.updatePosition).Methods("PUT")
	router.HandleFunc("/api/position/{id}", c.deletePosition).Methods("DELETE")
}

func (c *Controller) showPositionList(w http.ResponseWriter, r *http.Request) {

	ps, err := c.s.FindAllPositions()
	if err != nil {
		internalServerError(w, err)
		return
	}

	tData := struct {
		Positions []repository.Position
	}{
		ps,
	}
	err = c.t.ExecuteTemplate(w, "positionList", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) showNewPosition(w http.ResponseWriter, r *http.Request) {

	p := c.s.GetCurrentPosition()
	p.Name = time.Now().Format("2006-01-02 3:04:05 PM")

	tData := struct {
		Position repository.Position
	}{
		p,
	}
	err := c.t.ExecuteTemplate(w, "positionEditor", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) showExistingPosition(w http.ResponseWriter, r *http.Request) {

	id := getPathInt(r, "id")
	p, err := c.s.FindPosition(id)
	if err != nil {
		notFoundError(w, err)
		return
	}

	tData := struct {
		Position repository.Position
	}{
		p,
	}
	err = c.t.ExecuteTemplate(w, "positionEditor", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) createPosition(w http.ResponseWriter, r *http.Request) {

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

func (c *Controller) updatePosition(w http.ResponseWriter, r *http.Request) {

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

func (c *Controller) deletePosition(w http.ResponseWriter, r *http.Request) {
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
