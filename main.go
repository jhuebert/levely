package main

import (
	"context"
	"flag"
	"github.com/jhuebert/levely/controller"
	"github.com/jhuebert/levely/repository"
	"github.com/jhuebert/levely/service"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	var configPath string
	var dbPath string
	flag.StringVar(&configPath, "c", "", "path to the config file")
	flag.StringVar(&dbPath, "d", "levely.db", "path to the database file")
	flag.Parse()

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Error(err)
			flag.Usage()
			return
		}

		level, err := logrus.ParseLevel(viper.GetString(LogLevel))
		if err != nil {
			logrus.Warnf("Invalid input log level \"%v\". Can be one of trace, debug, info, warn, error, fatal, panic. Setting log level to info", viper.GetString(LogLevel))
			level = logrus.InfoLevel
		}
		logrus.SetLevel(level)
	}

	// Start a new router
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	r, err := repository.New(dbPath)
	if err != nil {
		logrus.Error(err)
		return
	}

	d := getDriver()
	s := service.New(r, d)
	rc := controller.New(s)
	rc.RegisterRoutes(router)

	// Define the REST server
	srv := &http.Server{
		Handler:      router,
		Addr:         viper.GetString(ServerAddress),
		WriteTimeout: viper.GetDuration(ServerWriteTimeout) * time.Millisecond,
		ReadTimeout:  viper.GetDuration(ServerReadTimeout) * time.Millisecond,
		IdleTimeout:  viper.GetDuration(ServerIdleTimeout) * time.Millisecond,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(ServerStopTimeout)*time.Millisecond)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	if err != nil {
		logrus.Error(err)
	}

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func getDriver() *i2c.MPU6050Driver {

	adaptor := raspi.NewAdaptor()
	err := adaptor.Connect()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	d := i2c.NewMPU6050Driver(adaptor, i2c.WithBus(viper.GetInt(DeviceI2CBus)), i2c.WithAddress(viper.GetInt(DeviceI2CAddress)))
	err = d.Start()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return d
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("%v %v", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
