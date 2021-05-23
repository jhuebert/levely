package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/jhuebert/levely/config"
	"github.com/jhuebert/levely/controller"
	"github.com/jhuebert/levely/repository"
	"github.com/jhuebert/levely/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	var configPath string
	var dbPath string
	flag.StringVar(&configPath, "c", "", "path to the config file")
	flag.StringVar(&dbPath, "d", "levely.db", "path to the database file")
	flag.Parse()

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logFormatter.FullTimestamp = true
	logrus.SetFormatter(logFormatter)

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Errorf("error reading config: %v", err)
			flag.Usage()
			return
		}

		level, err := logrus.ParseLevel(viper.GetString(config.LogLevel))
		if err != nil {
			logrus.Warnf("Invalid input log level \"%v\". Can be one of trace, debug, info, warn, error, fatal, panic. Setting log level to info", viper.GetString(config.LogLevel))
			level = logrus.InfoLevel
		}
		logrus.SetLevel(level)
	}

	// Start a new router
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	r, err := repository.New(dbPath)
	if err != nil {
		logrus.Errorf("could not open database: %v", err)
		return
	}

	//TODO Migrate DB

	d := getDriver()
	s := service.New(r, d)
	rc, err := controller.New(s)
	if err != nil {
		logrus.Errorf("error setting up controller: %v", err)
		return
	}
	rc.RegisterRoutes(router)

	// Define the REST server
	srv := &http.Server{
		Handler:      router,
		Addr:         viper.GetString(config.ServerAddress),
		WriteTimeout: viper.GetDuration(config.ServerWriteTimeout),
		ReadTimeout:  viper.GetDuration(config.ServerReadTimeout),
		IdleTimeout:  viper.GetDuration(config.ServerIdleTimeout),
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
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.ServerStopTimeout))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("closing database connection")
	r.Close()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Info("shutting down")
	os.Exit(0)
}

func getDriver() *i2c.MPU6050Driver {

	adaptor := raspi.NewAdaptor()
	err := adaptor.Connect()
	if err != nil {
		logrus.Errorf("could not connect to Raspberry Pi I2C bus: %v", err)
		return nil
	}

	d := i2c.NewMPU6050Driver(adaptor, i2c.WithBus(viper.GetInt(config.DeviceI2CBus)), i2c.WithAddress(viper.GetInt(config.AccelerometerI2CAddress)))
	err = d.Start()
	if err != nil {
		logrus.Errorf("could not communicate with accelerometer: %v", err)
		return nil
	}

	return d
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("%v %v", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
