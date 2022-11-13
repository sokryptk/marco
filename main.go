package main

import (
    "fmt"
    "me.kryptk.marco/services"
)

func main() {
	//
	w := services.NewNMWiFi()
	defer w.Conn.Close()

	devices := w.GetDevices()

	for _, r := range devices {
        connection, err := r.GetActiveConnection()
        if err != nil {
            continue
        }
        fmt.Println(connection.GetFrequency())
	}
}
