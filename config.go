package main

import "github.com/spf13/viper"

const (
	LogLevel           string = "log.level"
	ServerAddress      string = "server.address"
	ServerReadTimeout  string = "server.timeout.read"
	ServerWriteTimeout string = "server.timeout.write"
	ServerIdleTimeout  string = "server.timeout.idle"
	ServerStopTimeout  string = "server.timeout.stop"
	DeviceI2CBus       string = "device.i2c.bus"
	DeviceI2CAddress   string = "device.i2c.address"
)

func init() {
	viper.SetDefault(LogLevel, "info")
	viper.SetDefault(ServerAddress, ":8080")
	viper.SetDefault(ServerReadTimeout, "5000")
	viper.SetDefault(ServerWriteTimeout, "5000")
	viper.SetDefault(ServerIdleTimeout, "60000")
	viper.SetDefault(ServerStopTimeout, "15000")
	viper.SetDefault(DeviceI2CBus, "1")
	viper.SetDefault(DeviceI2CAddress, "0x68")
}
