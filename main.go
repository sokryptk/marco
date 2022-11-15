package main

import (
    "fmt"
    "me.kryptk.marco/repository"
    "me.kryptk.marco/services"
)

func main() {
    s := services.NewNM()
    devices := s.GetDevices()

    for _, d := range devices {
        switch d := d.(type) {
        case repository.WiFiDevice:
            err := d.RequestScan()
            if err != nil {
                return
            }

            ap := d.GetAccessPoints()
            for _, a := range ap {
                fmt.Println(a.GetSSID())
                fmt.Println(a.Connect())
            }

            fmt.Println(d.GetHwAddress())
        case repository.EthernetDevice:
            fmt.Println(d.GetHwAddress())
        }
    }
//    p := tea.NewProgram(models.NewHome())
//    if err := p.Start(); err != nil {
//        fmt.Printf("Alas, there's been an error: %v", err)
//        os.Exit(1)
//    }
}
