package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (c *Controller) registerCalibrationRoutes(router *mux.Router) {
	router.HandleFunc("/html/calibration", c.showCalibration).Methods("GET")
	router.HandleFunc("/api/calibration", c.updateCalibration).Methods("PUT")
}

func (c *Controller) showCalibration(w http.ResponseWriter, r *http.Request) {

	p := c.s.GetCalibration()

	tData := struct {
		Position repository.Position
	}{
		p,
	}
	err := c.t.ExecuteTemplate(w, "calibration", tData)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Controller) updateCalibration(w http.ResponseWriter, r *http.Request) {

	p, err := readPosition(r)
	if err != nil {
		internalServerError(w, err)
		return
	}

	p, err = c.s.UpdateCalibration(p)
	if err != nil {
		internalServerError(w, err)
		return
	}

	sendResponse(w, p)
}
