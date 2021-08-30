# Levely

# Introduction


# Hardware

## Minimum Hardware Requirements
- Raspberry Pi Zero W
- MPU-6050 connected via I2C
- Power source
    - USB
    - Direct 5V on GPIO

## Recommended Additions
- LM2591 Buck Converter

# System Configuration

# Wifi
- Wifi configuration as AP mode or not

# I2C
- Instructions to Enable i2c

### Increase I2C frequency to 400kHz
1. Open /boot/config.txt file
    - `sudo nano /boot/config.txt`
1. Modify the line containing `dtparam=i2c_arm=on`
    - `dtparam=i2c_arm=on,i2c_arm_baudrate=400000`
1. Reboot Raspberry Pi

## Power Reduction

### Turn off HDMI
```sh
sudo /opt/vc/bin/tvservice -o
```

### Turn off LED
```sh
echo none | sudo tee /sys/class/leds/led0/trigger
echo 1 | sudo tee /sys/class/leds/led0/brightness
```

# Schematic

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

### Values
- `passthrough`
    - This filter does not change the data coming from the accelerometer and simply returns the raw values. This might be useful for troubleshooting. The level data shown might be noisy when using this filter.
- `smoother`
    - This filter . This filter can be tuned by setting the `accelerometer.filter.smoother.smoothing` configuration property.
- `average` (default)
    - This filter performs a moving average on the data returned from the accelerometer. The size of the moving average filter is configured via the `accelerometer.filter.average.size` configuration property.

## accelerometer.filter.smoother.smoothing

## accelerometer.filter.average.size
This property configures the size of the moving average filter. This property is only applicable when `accelerometer.filter.selected=average`. Larger values will use more samples in the moving average resulting in a smoother level display at the expense of a slower response. Smaller values will have a faster response with the result of more noise in the displayed level value. The default value is `1000` which represents about two seconds at the default 500 samples/second.

### Example
`accelerometer.filter.average.size=1000`
