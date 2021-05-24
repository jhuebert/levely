package controller

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed static
var staticFiles embed.FS

func (c *Controller) registerStaticRoutes(router *mux.Router) {
	router.HandleFunc("/", redirectToHome).Methods("GET")
	cacheTime := viper.GetDuration(config.ServerCachePeriod)
	router.PathPrefix("/static/").Handler(cacheControlWrapper(cacheTime, http.FileServer(http.FS(staticFiles))))
}

func cacheControlWrapper(cacheTime time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int(cacheTime.Seconds())))
		h.ServeHTTP(w, r)
	})
}

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("redirecting to home")
	http.Redirect(w, r, "http://"+r.Host+"/home", http.StatusFound)
}
