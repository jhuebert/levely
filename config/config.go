package config

import "github.com/spf13/viper"

const (
	LogLevel string = "log.level"

	ServerAddress      string = "server.address"
	ServerReadTimeout  string = "server.timeout.read"
	ServerWriteTimeout string = "server.timeout.write"
	ServerIdleTimeout  string = "server.timeout.idle"
	ServerStopTimeout  string = "server.timeout.stop"
	ServerCachePeriod  string = "server.cache.period"

	DeviceI2CBus string = "device.i2c.bus"

	DisplayLevelTolerance string = "display.level.tolerance"
	DisplayUpdateRate     string = "display.update.rate"

	AccelerometerI2CAddress        string = "accelerometer.i2c.address"
	AccelerometerUpdateSleepWait   string = "accelerometer.update.sleep.wait"
	AccelerometerUpdateSleepPeriod string = "accelerometer.update.sleep.period"
	AccelerometerUpdatePeriod      string = "accelerometer.update.period"

	AccelerometerFilterSelected          string = "accelerometer.filter.selected"
	AccelerometerFilterSmootherSmoothing string = "accelerometer.filter.smoother.smoothing"
	AccelerometerFilterAverageSize       string = "accelerometer.filter.average.size"

	FilterSmoother    string = "smoother"
	FilterAverage     string = "average"
	FilterPassthrough string = "passthrough"
)

func init() {
	viper.SetDefault(LogLevel, "info")

	// set web server defaults
	viper.SetDefault(ServerAddress, ":8080")
	viper.SetDefault(ServerReadTimeout, "5s")
	viper.SetDefault(ServerWriteTimeout, "5s")
	viper.SetDefault(ServerIdleTimeout, "60s")
	viper.SetDefault(ServerStopTimeout, "15s")
	viper.SetDefault(ServerCachePeriod, "1h")

	// set device defaults
	viper.SetDefault(DeviceI2CBus, "1")

	// set display defaults
	viper.SetDefault(DisplayLevelTolerance, "0.1")
	viper.SetDefault(DisplayUpdateRate, "4")

	// set accelerometer defaults
	viper.SetDefault(AccelerometerI2CAddress, "0x68")
	viper.SetDefault(AccelerometerUpdateSleepWait, "15s")
	viper.SetDefault(AccelerometerUpdateSleepPeriod, "500ms")
	viper.SetDefault(AccelerometerUpdatePeriod, "5ms")

	// set accelerometer filter defaults
	viper.SetDefault(AccelerometerFilterSelected, "average")
	viper.SetDefault(AccelerometerFilterSmootherSmoothing, "1000")
	viper.SetDefault(AccelerometerFilterAverageSize, "400")
}
