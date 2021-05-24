package config

import "github.com/spf13/viper"

const (
	LogLevel string = "log.level"

	ServerAddress                  string = "server.address"
	ServerReadTimeout              string = "server.timeout.read"
	ServerWriteTimeout             string = "server.timeout.write"
	ServerIdleTimeout              string = "server.timeout.idle"
	ServerStopTimeout              string = "server.timeout.stop"
	ServerStaticCacheControlPeriod string = "server.static.cache.control.period"

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
	AccelerometerFilterMedianSize        string = "accelerometer.filter.median.size"
	AccelerometerFilterIirPath           string = "accelerometer.filter.iir.path"
	AccelerometerFilterIirPreset         string = "accelerometer.filter.iir.preset"
	AccelerometerFilterFirPath           string = "accelerometer.filter.fir.path"
	AccelerometerFilterFirPreset         string = "accelerometer.filter.fir.preset"

	FilterSmoother    string = "smoother"
	FilterAverage     string = "average"
	FilterMedian      string = "median"
	FilterIir         string = "iir"
	FilterFir         string = "fir"
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
	viper.SetDefault(ServerStaticCacheControlPeriod, "1h")

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
	viper.SetDefault(AccelerometerFilterMedianSize, "200")
	viper.SetDefault(AccelerometerFilterMedianSize, "200")
	viper.SetDefault(AccelerometerFilterIirPath, "")
	viper.SetDefault(AccelerometerFilterIirPreset, "1")
	viper.SetDefault(AccelerometerFilterFirPath, "")
	viper.SetDefault(AccelerometerFilterFirPreset, "1")
}
