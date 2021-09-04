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

## General

### log.level
Specifies the level of logging that the program should output to standard out. The default log level is `info`. This will normally not need to be changed, but setting the log level to `debug` can assist in troubleshooting issues.

#### Values
Value selected will output log statements at the level selected and those above it in the list.
- `panic`
- `fatal`
- `error`
- `warn`
- `info` (default)
- `debug`
- `trace`

## Server

### `server.address`
This property defines the address where the application can be reached from the client. The default value is `:8080`.

### `server.timeout.read`
This property defines the maximum length of time the application will wait for the incoming client connection to send its data. The default value is `5s`.

### `server.timeout.write`
This property defines the maximum length of time the application will allow the outgoing client connection to be written to. This value also affects how long a server sent event (SSE) connecction can be open before a new connection will have to be made. The default value is `60s`.

### `server.timeout.idle`
This property defines the maximum length of time that a client connection can be idle before it is closed. The default value is `60s`.

### `server.timeout.stop`
This property defines the maximum length of time the application will wait for the web server to stop. The default value is `15s`.

### `server.cache.period`
This property defines the length of time client static assets should be cached before refetching. The default value is `1h`.

## Device

### `device.i2c.bus`
This property defines the identifier of the I2C bus on the host device. The default value is `1`.

## Display

### `display.level.tolerance`
This property defines the angle in degrees that below which the client will indicate level has been achieved. The default value is `0.1` degrees.

### `display.update.rate`
This property defines how frequenty the client display should be updated with the latest level data. The default value is `4` updates/second.

### `display.sse.enabled`
This property defines whether or not server sent events (SSE) support is enabled on the client. Polling would be used if SSE is disabled. SSE is more efficient than polling. The default value is `true`.

## Accelerometer

### `accelerometer.i2c.address`
This property defines the I2C address of the accelerometer. This would need to be changed if the default address of the accelerometer is changed on the device. The default value is `0x68`.

### `accelerometer.update.sleep.wait`
This property defines how long the application will wait with no clients requesting samples before it will go to sleep. The default value is `2s`.

### `accelerometer.update.sleep.period`
This property defines how often the accelerometer is polled when the application is sleeping. Shorter periods result in more updated values when first viewing the current level at the cost of higher energy usage. Longer values also result in a longer time to the first value displayed when viewing the current level. The default value is `500ms`.

### `accelerometer.update.period`
This property defines how frequently the accelerometer is polled when clients are actively requesting level data. The minimum value for this property is 1ms which represents the maximum update rate of the accelerometer. Smaller values result in more data that can be filtered at the expense of higher CPU usage. The default value is `2ms` which represents 500 samples/second.

## Filter

### `accelerometer.filter.selected`
This property defines which method should be used to filter the raw accelerometer values.

#### Values
- `passthrough`
    - This filter does not change the data coming from the accelerometer and simply returns the raw values. This might be useful for troubleshooting. The level data shown might be noisy when using this filter.
- `smoother`
    - This filter works similarly to a moving average filter. It is fast to calculate and requires no additional memory for large smoothing values. It works by using the difference in value between last and current samples and only adding a portion of the difference to the previous output value. This filter can be tuned by setting the `accelerometer.filter.smoother.smoothing` configuration property.
- `average` (default)
    - This filter performs a moving average on the data returned from the accelerometer. The size of the moving average filter is configured via the `accelerometer.filter.average.size` configuration property.

### `accelerometer.filter.smoother.smoothing`
This property configures the smoother filter. Larger values result in smoother output at the expense of temporal response. The default value is `250` which represents about a half-second response time using the default 500 samples/second. 

### `accelerometer.filter.average.size`
This property configures the size of the moving average filter. This property is only applicable when `accelerometer.filter.selected=average`. Larger values will use more samples in the moving average resulting in a smoother level display at the expense of a slower response. Smaller values will have a faster response with the result of more noise in the displayed level value. Changing the sampling rate will require the filter size to be changed in proportion. The default value is `1000` which represents about a two seconds response time at the default 500 samples/second.
