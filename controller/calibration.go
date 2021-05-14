package controller

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) registerCalibrationRoutes(router *mux.Router) {
	router.HandleFunc("/api/calibration", c.GetCalibration).Methods("GET")
	router.HandleFunc("/api/calibration", c.UpdateCalibration).Methods("PUT")
}

func (c *Controller) GetCalibration(w http.ResponseWriter, r *http.Request) {
	p := c.s.GetCalibration()
	sendResponse(w, p)
}

func (c *Controller) UpdateCalibration(w http.ResponseWriter, r *http.Request) {

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
