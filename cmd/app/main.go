package main

import (
    "fmt"
    "github.com/dantheman213/go-cli"
    "github.com/dantheman213/gps"
    "strings"
    "time"
)

func main() {
    gps := gps.NewGPS()
    // TODO: ingest NMEA sentences from serial device

    for true {
        loc := gps.GetGPSLocationInDMSPretty()
        fmt.Println(loc)

        _, output, _, err := cli.MakeAndRunCommandThenWait("iwlist wlan0 scanning | egrep 'Cell |Encryption|Quality|Last beacon|ESSID'")
        if err != nil {
            panic(err)
        }

        for _, line := range strings.Split(output.String(), "\n") {
            fmt.Println(line)
        }

        time.Sleep(time.Duration(3) * time.Second)
    }
}
