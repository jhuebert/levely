package controller

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func (c *Controller) registerStaticRoutes(router *mux.Router) {
	router.HandleFunc("/", getFile).Methods("GET")
	router.HandleFunc("/favicon.ico", getFile).Methods("GET")
	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFiles)))
}

func getFile(w http.ResponseWriter, r *http.Request) {

	fp := r.RequestURI[1:]
	if fp == "" {
		fp = "static/index.html"
	} else if fp == "favicon.ico" {
		fp = "static/favicon.ico"
	}

	data, err := staticFiles.ReadFile(fp)
	if err != nil {
		notFoundError(w, err)
	}

	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		logrus.Error(err)
	}
}
