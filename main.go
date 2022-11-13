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
        fmt.Println(r.GetHwAddresss(), r.GetDeviceType())
        
        err := r.RequestScan()
        if err != nil {
            log.Println(err)
            continue

        }

		points := r.GetAccessPoints()

		for _, r := range points {
            fmt.Println(r.GetSSID(), fmt.Sprintf("%.1fGHz", float64(r.GetFrequency()) / float64(1000)), r.GetStrength())
		}
	}
}
