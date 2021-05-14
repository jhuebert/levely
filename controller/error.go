package controller

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func internalServerError(w http.ResponseWriter, err error) {
	logrus.Error(err)
	sendError(w, err, http.StatusInternalServerError)
}

//TODO Take identifier as input
func notFoundError(w http.ResponseWriter, err error) {
	logrus.Error(err)
	sendError(w, err, http.StatusNotFound)
}

func sendError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
