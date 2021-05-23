package controller

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/service"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	s *service.Service
	t *template.Template
}

//go:embed html
var templateFiles embed.FS

func New(s *service.Service) (Controller, error) {
	t, err := template.ParseFS(templateFiles, "html/*.gohtml")
	return Controller{s, t}, err
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	c.registerStaticRoutes(router)
	c.registerPositionRoutes(router)
	c.registerDeviceRoutes(router)
	c.registerCalibrationRoutes(router)
	c.registerPreferenceRoutes(router)
	c.registerHomeRoutes(router)
	c.registerRecallRoutes(router)
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
