# Levely

# Introduction


# Hardware

## Minimum Hardware Requirements
- Raspberry Pi Zero W
- MPU-6050 connected via I2C

## Recommended Additions
- LM2591 Buck Converter

# Running
Usage of levely:
  -c string
        path to the config file
  -d string
        path to the database file (default "levely.db")


# Configuration
Designed to be run without any modifications, but some installations might want to customize some of the properties. Sane defaults were selected for all properties and only those that need to be modified should be changed. If a default value is not working for most people, we should change it.

Configuration can stored in JSON, TOML, YAML, HCL, envfile or Java properties file.

## log.level
Specifies the level of logging that the program should output to standard out. The default log level is `info`. This will normally not need to be changed, but setting the log level to `debug` can assist in troubleshooting issues.

### Example
`log.level=debug`

### Values
Value selected will output log statements at the level selected and those above it in the list.
- panic
- fatal
- error
- warn
- info (default)
- debug
- trace

## server.address

## server.timeout.read

## server.timeout.write

## server.timeout.stop

## server.cache.period

## device.i2c.bus

## display.level.tolerance

## display.update.rate

## display.sse.enabled

## accelerometer.i2c.address

## accelerometer.update.sleep.wait

## accelerometer.update.sleep.period

## accelerometer.update.period

## accelerometer.filter.selected

## accelerometer.filter.smoother.smoothing

## accelerometer.filter.average.size

