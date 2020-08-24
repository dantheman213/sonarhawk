# sonarhawk

Create precise maps of WiFi networks using commodity GPS hardware and a portable computer. Supports Linux, MacOS, and Windows.

## Getting Started

You'll need a USB GPS receiver using serial over USB connection. u-blox 7 or similar will work great and is inexpensive. This utility will use the default WiFi device to scan for networks. 

```
git clone https://github.com/dantheman213/sonarhawk
cd sonarhawk/
make deps
make # NOTE: run 'make windows' or 'make macos' if not on linux
cd bin/

# You'll want to walk or drive around while running this utility to sample the area that you want. If driving recommended to go slow so device has time to scan networks in the area.
# NOTE: Add ".exe" after binary if on Windows
./sonarhawk-survey -output survey.csv

# This can be run at any time after the survey has been completed.
./sonarkhawk-process -input survey.csv -output survey.kml
```

Open `survey.kml` in **Google Earth Pro** for a detailed map of the WiFi networks that you surveyed.

## Hardware

### GPS

Any GPS dongle using USB-over-Serial should work. Confirmed working hardware items are:

* BN-80U

* BS-708

* BU-353-S4

* VK-162

* VK-172
