package main

import (
    "fmt"
    "log"
    "me.kryptk.marco/services"
)

func main() {
	//
	w := services.NewNMWiFi()
	defer w.Conn.Close()

	devices := w.GetDevices()

	for _, r := range devices {
        err := r.RequestScan()
        if err != nil {
            log.Println(err)
            continue

        }

        fmt.Println(r.GetHwAddresss())
		points := r.GetAccessPoints()

		for _, r := range points {
            fmt.Println(r.GetSSID(), fmt.Sprintf("%.1fGHz", float64(r.GetFrequency()) / float64(1000)), r.GetStrength())
		}
	}
}
