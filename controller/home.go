package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (c *Controller) registerHomeRoutes(router *mux.Router) {
	router.HandleFunc("/html/home", c.showHome).Methods("GET")
}

func (c *Controller) showHome(w http.ResponseWriter, r *http.Request) {

	ps, err := c.s.FindFavoritePositions()
	if err != nil {
		internalServerError(w, err)
		return
	}

	tData := struct {
		Positions []repository.Position
	}{
		ps,
	}
	err = c.t.ExecuteTemplate(w, "home", tData)
	if err != nil {
		logrus.Error(err)
	}
}
