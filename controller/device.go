package controller

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) registerDeviceRoutes(router *mux.Router) {
	router.HandleFunc("/api/corrected", c.GetCurrentPosition).Methods("GET")
	router.HandleFunc("/api/uncorrected", c.GetUncorrected).Methods("GET")
}

func (c *Controller) GetCurrentPosition(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetCurrentPosition())
}

func (c *Controller) GetUncorrected(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetUncorrected())
}
