package main

import (
    "fmt"
    "me.kryptk.marco/repository"
    "me.kryptk.marco/services"
)

func main() {
	//
	var w repository.WiFi =  services.NewNMWiFi()
	defer w.Close()

	devices := w.GetDevices()

	for _, r := range devices {
        connection, err := r.GetActiveConnection()
        if err != nil {
            continue
        }
        fmt.Println(connection.GetFrequency())
	}
}
