package controller

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func internalServerError(w http.ResponseWriter, err error) {
	sendError(w, err, http.StatusInternalServerError)
}

func notFoundError(w http.ResponseWriter, err error) {
	sendError(w, err, http.StatusNotFound)
}

func sendError(w http.ResponseWriter, err error, status int) {
	logrus.Error(err)
	http.Error(w, err.Error(), status)
}
