package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (c *Controller) registerDeviceRoutes(router *mux.Router) {
	router.HandleFunc("/api/corrected", c.getCurrentPosition).Methods("GET")
	router.HandleFunc("/api/uncorrected", c.getUncorrected).Methods("GET")
	router.HandleFunc("/api/corrected/event", c.getCurrentPositionEventStream).Methods("GET")
}

func (c *Controller) getCurrentPosition(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetCurrentPosition())
}

func (c *Controller) getUncorrected(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, c.s.GetUncorrected())
}

func (c *Controller) getCurrentPositionEventStream(w http.ResponseWriter, r *http.Request) {

	// grab the time this event stream started so that we can close the connection before the write timeout expires
	start := time.Now()

	// ensure that the writer supports flushing so that events can be sent immediately to the client
	flusher, ok := w.(http.Flusher)
	if !ok {
		internalServerError(w, errors.New("streaming is not supported"))
		return
	}

	// set the response headers so that the client knows that the response is an event stream
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// indicate to the client that it should immediately open a new event stream after this one closes
	err := sendRetry(w, 0)
	if err != nil {
		logrus.Errorf("error sending retry: %v", err)
		return
	}

	// send position updates to the client based on the configured rate
	updateDuration := (1 * time.Second) / time.Duration(viper.GetFloat64(config.DisplayUpdateRate))

	// only write to this event stream until right before the server write timeout so that the connection is closed gracefully
	writeTimeout := viper.GetDuration(config.ServerWriteTimeout) - updateDuration

	for range time.Tick(updateDuration) {

		// don't send any more events that might result in a timeout being hit
		if time.Now().Sub(start) > writeTimeout {
			return
		}

		p := c.s.GetCurrentPosition()

		// check that the connection is still open right before attempting to send
		if isClosed(r) {
			logrus.Debug("event stream closed")
			return
		}

		err = sendData(w, p)
		if err != nil {
			logrus.Errorf("error sending data: %v", err)
			return
		}

		// ensure that the event is sent to the client immediately
		flusher.Flush()
	}
}

func isClosed(r *http.Request) bool {
	select {
	case <-r.Context().Done():
		logrus.Debug("connection closed")
		return true
	default:
		logrus.Debug("connection open")
		return false
	}
}

func sendRetry(w http.ResponseWriter, delay int) error {
	logrus.Debug("sending retry event: %d", delay)
	_, err := fmt.Fprintf(w, "retry: %d\n\n", delay)
	return err
}

func sendData(w http.ResponseWriter, d interface{}) error {
	logrus.Debug("sending data event: %v", d)
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", j)
	return err
}
