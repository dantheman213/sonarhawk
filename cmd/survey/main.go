package main

import (
    "flag"
    "fmt"
    "github.com/dantheman213/go-cli"
    libGPS "github.com/dantheman213/go-gps"
    "github.com/dantheman213/sonarhawk/pkg/ingest"
    log "github.com/sirupsen/logrus"
    "os"
    "regexp"
    "runtime"
    "strconv"
    "strings"
    "time"
)

// netsh wlan show networks mode=bssid
func main() {
    output := flag.String("output", "", "The file path to write the output at.")
    flag.Parse()

    if *output == "" {
        log.Fatal("no file path for output selected")
    }

    gps := libGPS.NewGPS()
    go ingest.IngestGPSData(gps)

    f, err := os.Create(*output)
    if err != nil {
        log.Fatal(err)
    }

    header := "SSID, Authentication, Encryption, BSSID, Radio Type, Signal, Latitude, Longitude\n"
    if _, err := f.WriteString(header); err != nil {
        log.Fatal(err)
    }

    for true {
        time.Sleep(2 * time.Second)

        _, out, _, err := cli.MakeAndRunCommandThenWait(generateWifiCommand())
        if err != nil {
            log.Fatal(err)
        }

        wifi, err := ingestWifiData(out.String())
        if err != nil {
            log.Error(err)
            continue
        }

        loc, err := gps.GetGPSLocation()
        if err != nil {
            log.Error(err)
            continue
        }

        point := &DataPoint{
            WiFi:      *wifi,
            Latitude:  loc.Latitude,
            Longitude: loc.Longitude,
        }

        csv := fmt.Sprintf("%s, %s, %s, %s, %s, %v, %v, %v\n", point.WiFi.SSID, point.WiFi.Authentication, point.WiFi.Encryption, point.WiFi.BSSID, point.WiFi.RadioType, point.WiFi.Signal, point.Latitude, point.Longitude)
        log.Info(csv)
        if _, err := f.WriteString(csv); err != nil {
            log.Error(err)
        }
    }

    _ = f.Close()
}

func generateWifiCommand() string {
    if runtime.GOOS == "windows" {
        return "netsh wlan show networks mode=bssid"
    }

    return ""
}

type DataPoint struct {
    WiFi WiFiData
    Latitude float64
    Longitude float64
}

type WiFiData struct {
    SSID string
    Authentication string
    Encryption string
    BSSID string
    Signal float64
    RadioType string
}

func ingestWifiData(dat string) (*WiFiData, error) {
    if runtime.GOOS == "windows" {
        return ingestWifiDataWindows(dat)
    }

    return nil, nil
}

// example payload
//SSID 1 : MyWiFi-5G
//Network type            : Infrastructure
//Authentication          : WPA2-Personal
//Encryption              : CCMP
//BSSID 1                 : 92:1e:19:5b:3d:47
//Signal             : 15%
//Radio type         : 802.11ac
//Channel            : 44
//Basic rates (Mbps) : 6 12 24
//Other rates (Mbps) : 9 18 36 48 54
func ingestWifiDataWindows(dat string) (*WiFiData, error) {
    dat = strings.ReplaceAll(dat, "\r", "")
    lines := strings.Split(dat, "\n")

    result := &WiFiData{}
    for _, line := range lines {
        line = strings.TrimSpace(line)

        if strings.HasPrefix(line, "SSID") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")
            result.SSID = strings.TrimSpace(parts[1])
        } else if strings.HasPrefix(line, "Authentication") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")
            result.Authentication = strings.TrimSpace(parts[1])
        } else if strings.HasPrefix(line, "Encryption") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")
            result.Encryption = strings.TrimSpace(parts[1])
        } else if result.BSSID == "" && strings.HasPrefix(line, "BSSID") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")
            result.BSSID = strings.TrimSpace(parts[1])
        } else if  result.Signal == 0 && strings.HasPrefix(line, "Signal") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")

            reg, err := regexp.Compile("[^0-9]+")
            if err != nil {
                log.Error(err)
            }

            part := reg.ReplaceAllString(parts[1], "")
            num, err := strconv.ParseFloat(part, 32)
            if err != nil {
                return nil, err
            }
            result.Signal = num / 100
        } else if result.RadioType == "" && strings.HasPrefix(line, "Radio type") && strings.Index(line, ":") > -1 {
            parts := strings.Split(line, ":")
            result.RadioType = strings.TrimSpace(parts[1])
        }
    }

    return result, nil
}