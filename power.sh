# Increase I2C frequency to 400kHz

# Turn off HDMI
sudo /opt/vc/bin/tvservice -o

# Turn off LED
echo none | sudo tee /sys/class/leds/led0/trigger
echo 1 | sudo tee /sys/class/leds/led0/brightness
