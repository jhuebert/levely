package controller

import (
	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/config"
	"github.com/spf13/viper"
	"net/http"
)

func (c *Controller) registerConfigRoutes(router *mux.Router) {
	router.HandleFunc("/api/config", c.GetConfig).Methods("GET")
}

type Config struct {
	DisplayLevelTolerance float64 `json:"displayLevelTolerance"`
	DisplayUpdateRate     float64 `json:"displayUpdateRate"`
}

func (c *Controller) GetConfig(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, Config{
		DisplayLevelTolerance: viper.GetFloat64(config.DisplayLevelTolerance),
		DisplayUpdateRate:     viper.GetFloat64(config.DisplayUpdateRate),
	})
}
