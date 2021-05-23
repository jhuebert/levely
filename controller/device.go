package controller

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) registerDeviceRoutes(router *mux.Router) {
	router.HandleFunc("/api/corrected", c.getCurrentPosition).Methods("GET")
	router.HandleFunc("/api/uncorrected", c.getUncorrected).Methods("GET")
}

func (c *Controller) getCurrentPosition(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetCurrentPosition())
}

func (c *Controller) getUncorrected(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetUncorrected())
}
