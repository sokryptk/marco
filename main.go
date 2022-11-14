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
