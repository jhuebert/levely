package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Controller struct {
	s *service.Service
}

func New(s *service.Service) Controller {
	return Controller{s}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	c.registerStaticRoutes(router)
	c.registerPositionRoutes(router)
	c.registerDeviceRoutes(router)
	c.registerCalibrationRoutes(router)
	c.registerPreferenceRoutes(router)
}

func getPathInt(r *http.Request, name string) int {
	params := mux.Vars(r)
	value, _ := strconv.Atoi(params[name])
	logrus.Debugf("%v: %v", name, value)
	return value
}

func sendResponse(w http.ResponseWriter, v interface{}) {

	js, err := json.Marshal(v)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		logrus.Error(err)
	}
}
