package ingest

import (
    libGPS "github.com/dantheman213/go-gps"
    "github.com/dantheman213/go-gps/serial"
    log "github.com/sirupsen/logrus"
)

func IngestGPSData(engine *libGPS.GPS) {
    device, err := serial.DetectGPSDevice()
    if err != nil {
        log.Fatalf("[error] %s", err)
    }

    log.Info("connected to GPS device")

    for true {
        dat, err := serial.ReadSerialData(device.Port)
        if err != nil {
            log.Printf("couldn't read data stream")
            return
        }

        engine.IngestNMEASentences(dat)
    }
}

