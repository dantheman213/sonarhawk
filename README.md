# sonarhawk

Create precise maps of WiFi networks using commodity GPS hardware and a portable computer. Supports Linux, MacOS, and Windows.

## Getting Started

You'll need a USB GPS receiver using serial over USB connection. u-blox 7 or similar will work great and is inexpensive. This utility will use the default WiFi device to scan for networks. 

```
git clone https://github.com/dantheman213/sonarhawk
cd sonarhawk/
make deps
make
cd bin/
./sonarhawk-survey -output survey.csv
./sonarkhawk-process -input survey.csv -output survey.kml
```

Open `survey.kml` in Google Earth Pro for a detailed map of the WiFi networks that you surveyed.
