package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "github.com/dantheman213/sonarhawk/pkg/ingest"
    "github.com/dantheman213/sonarhawk/pkg/kml"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "math"
    "sort"
    "strconv"
    "strings"
)

func main () {
    input := flag.String("input", "", "The file path to read the CSV produced from survey binary.")
    output := flag.String("output", "", "The file path to write the KML output at.")
    flag.Parse()

    if *input == "" {
        log.Fatal("no file path for input selected")
    }

    if *output == "" {
        log.Fatal("no file path for output selected")
    }

    items, err := parseCSV(*input)
    if err != nil {
        log.Fatal(err)
    }

    kmlStr := compute(items)
    err = ioutil.WriteFile(*output, []byte(kmlStr), 0644)
    if err != nil {
        log.Fatal(err)
    }

    log.Info("COMPLETE!")
}

func parseCSV(path string) (map[string][]ingest.DataPoint, error) {
    results := make(map[string][]ingest.DataPoint, 0)

    dat, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    r := csv.NewReader(strings.NewReader(string(dat)))
    records, err := r.ReadAll()
    if err != nil {
        return nil, err
    }

    for i, record := range records {
        if i == 0 {
            continue
        }

        latitude, err := strconv.ParseFloat(record[6], 32)
        if err != nil {
            log.Error(err)
            continue
        }

        longitude, err := strconv.ParseFloat(record[7], 32)
        if err != nil {
            log.Error(err)
            continue
        }

        signal, err := strconv.ParseFloat(record[5], 32)
        if err != nil {
            log.Error(err)
            continue
        }

        result := ingest.DataPoint{
            Latitude:  latitude,
            Longitude: longitude,
            Wifi:      &ingest.WiFiData{
                SSID:           record[0],
                Authentication: record[1],
                Encryption:     record[2],
                BSSID:          strings.ToLower(record[3]),
                Signal:         signal,
                RadioType:      record[4],
            },
        }

        if results[result.Wifi.BSSID] == nil {
            results[result.Wifi.BSSID] = make([]ingest.DataPoint, 0)
        }
        results[result.Wifi.BSSID] = append(results[result.Wifi.BSSID], result)
    }

    return results, nil
}

func calculateGPSCenter(list *[]ingest.DataPoint) (float64, float64) {
    // convert to 3d cartesian from 2d
    x := 0.0
    y := 0.0
    z := 0.0

    for _, point := range *list {
        lat := point.Latitude * math.Pi / 180
        lng := point.Longitude * math.Pi / 180

        x += math.Cos(lat) * math.Cos(lng)
        y += math.Cos(lat) * math.Sin(lng)
        z += math.Sin(lat)
    }

    total := float64(len(*list))
    x = x / total
    y = y / total
    z = z / total

    targetLng := math.Atan2(y, x)
    targetSqrt := math.Sqrt(x * x + y * y)
    targetLat := math.Atan2(z, targetSqrt)

    // convert from 3d to 2d and return
    lat := targetLat * 180 / math.Pi
    lng := targetLng * 180 / math.Pi

    return lat, lng
}

func compute(items map[string][]ingest.DataPoint) string {
    payload := kml.TemplateParent
    placemarksStr := ""

    centerLat := 0.0
    centerLng := 0.0

    for _, item := range items {
        sort.Slice(item, func(i, j int) bool {
            // helps determine best signal rate captured and may be used in weighting for future GPS center calculation
            return item[i].Wifi.Signal > item[j].Wifi.Signal
        })

        if centerLat == 0 || centerLng == 0 {
            centerLat = item[0].Latitude
            centerLng = item[0].Longitude
        }

        lat, lng := calculateGPSCenter(&item)

        str := kml.TemplatePlacemark
        desc := fmt.Sprintf("%0.2f%% | %s | %s| %s | %s", item[0].Wifi.Signal * 100, item[0].Wifi.BSSID, item[0].Wifi.Authentication, item[0].Wifi.Encryption, item[0].Wifi.RadioType)
        str = strings.Replace(str, "%%KML_PLACEMARK_TITLE%%", xmlEscapeString(item[0].Wifi.SSID), 1)
        str = strings.Replace(str, "%%KML_PLACEMARK_DESCRIPTION%%", desc, 1)
        str = strings.Replace(str, "%%KML_PLACEMARK_LONGITUDE%%", fmt.Sprintf("%f", lng), 1)
        str = strings.Replace(str, "%%KML_PLACEMARK_LATITUDE%%", fmt.Sprintf("%f", lat), 1)

        placemarksStr += str + "\n"
    }

    payload = strings.Replace(payload, "%%KML_TITLE%%", "WiFi Site Survey", 1)
    payload = strings.Replace(payload, "%%KML_DESCRIPTION%%", fmt.Sprintf("Total networks surveyed: %d", len(items)), 1)

    payload = strings.Replace(payload, "%%KML_LOOKAT_LATITUDE%%", fmt.Sprintf("%f", centerLat), 1)
    payload = strings.Replace(payload, "%%KML_LOOKAT_LONGITUDE%%", fmt.Sprintf("%f", centerLng), 1)

    payload = strings.Replace(payload, "%%KML_PLACEMARKS%%", placemarksStr, 1)

    return payload
}

func xmlEscapeString(str string) string {
    str = strings.ReplaceAll(str, "\"", "&quot;")
    str = strings.ReplaceAll(str, "'", "&apos;")
    str = strings.ReplaceAll(str, "&", "&amp;")
    str = strings.ReplaceAll(str, "<", "&lt;")
    str = strings.ReplaceAll(str, ">", "&gt;")
    return str
}
