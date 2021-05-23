package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"github.com/jhuebert/levely/service"
	"github.com/sirupsen/logrus"
)

func (c *Controller) registerRecallRoutes(router *mux.Router) {
	router.HandleFunc("/html/level", c.showLevelRecall).Methods("GET")
	router.HandleFunc("/html/position/{id}/recall", c.showPositionRecall).Methods("GET")
}

func (c *Controller) showPositionRecall(w http.ResponseWriter, r *http.Request) {

	id := getPathInt(r, "id")
	p, err := c.s.FindPosition(id)
	if err != nil {
		notFoundError(w, err)
		return
	}

	tData := struct {
		Position    *repository.Position
		Preferences repository.Preferences
		Config      service.Config
	}{
		&p,
		c.s.GetPreferences(),
		c.s.GetConfig(),
	}

	err = c.t.ExecuteTemplate(w, "recall", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) showLevelRecall(w http.ResponseWriter, r *http.Request) {

	tData := struct {
		Position    *repository.Position
		Preferences repository.Preferences
		Config      service.Config
	}{
		nil,
		c.s.GetPreferences(),
		c.s.GetConfig(),
	}

	err := c.t.ExecuteTemplate(w, "recall", tData)
	if err != nil {
		logrus.Error(err)
	}
}
