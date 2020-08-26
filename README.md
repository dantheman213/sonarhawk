# sonarhawk

Create precise maps of WiFi networks using commodity GPS hardware and a portable computer. Supports Linux, MacOS, and Windows.

## Getting Started

You'll need a USB GPS receiver using serial over USB connection.Please see below for a list of compatible hardware. This utility will use the default WiFi device to scan for networks. 

### Setup
```
git clone https://github.com/dantheman213/sonarhawk
cd sonarhawk/
make deps
make # NOTE: run 'make windows' or 'make macos' if not on linux
cd bin/
```

### Survey
```
# You'll want to walk or drive around while running this utility to sample the area that you want. If driving recommended to go slow so device has time to scan networks in the area.
# NOTE: Add ".exe" after binary if on Windows
./sonarhawk-survey -output survey.csv
```

### Process Data

```
# This can be run at any time after the survey has been completed.
./sonarkhawk-process -input survey.csv -output survey.kml
```

### View Survey Results in GUI

Open `survey.kml` in [Google Earth Pro](https://www.google.com/earth/versions/#earth-pro) for a detailed map of the WiFi networks that you surveyed.

![Google Earth Pro Screenshot](https://raw.githubusercontent.com/dantheman213/sonarhawk/master/docs/example.jpg)

## Hardware

### GPS

Any GPS dongle using USB-over-Serial should work. Confirmed working hardware items are:

* BN-80U

* BS-708

* BU-353-S4

* VK-162

* VK-172
