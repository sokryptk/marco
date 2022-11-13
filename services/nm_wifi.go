package services

import (
    "github.com/godbus/dbus/v5"
    "log"
)

const nmInterface = "org.freedesktop.NetworkManager"

type NMWiFi struct {
    Conn *dbus.Conn
}

func NewNMWiFi() NMWiFi {
    conn, err := dbus.ConnectSystemBus()

    if err != nil {
        log.Fatal(err)
        return NMWiFi{}
    }

    return NMWiFi{
        Conn: conn,
    }
}

type NMDevice struct {
    conn *dbus.Conn
    Path dbus.ObjectPath
}

type NMAccessPoint struct {
    conn *dbus.Conn
    Path dbus.ObjectPath
}

func (ap NMAccessPoint) GetSSID() (s string) {
    v, err := ap.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
    if err != nil {
        log.Println(err)
    }

    return string(v.Value().([]byte))
}

func (dev NMDevice) RequestScan() error {
    err := dev.conn.Object(nmInterface, dev.Path).Call("org.freedesktop.NetworkManager.Device.Wireless.RequestScan", 0, map[string]interface{}{})
    if err != nil {
        return err.Err
    }

    return nil
}

func (n NMWiFi) GetDevices() []NMDevice {
    var devicePaths []dbus.ObjectPath
    err := n.Conn.Object(nmInterface, "/org/freedesktop/NetworkManager").Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devicePaths)
    if err != nil {
        return nil
    }

    devices := make([]NMDevice, len(devicePaths))

    for i, d := range devicePaths {
        devices[i] = NMDevice{conn: n.Conn, Path: d}
    }

    return devices
}

func (ap NMAccessPoint) GetFrequency() uint {
    v, err := ap.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Frequency")
    if err != nil {
        log.Println(err)
    }

    var frequency uint
    err = v.Store(&frequency)
    if err != nil {
        log.Println(err)
        return 0
    }

    return frequency
}

func (ap NMAccessPoint) GetStrength() uint {
    v, err := ap.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Strength")
    if err != nil {
        log.Println(err)
    }

    var strength uint
    err = v.Store(&strength)
    if err != nil {
        log.Println(err)
        return 0
    }

    return strength
}

func (dev NMDevice) GetAccessPoints() []NMAccessPoint  {
    var accessPaths []dbus.ObjectPath
    err := dev.conn.Object(nmInterface, dev.Path).Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&accessPaths)
    if err != nil {
        return nil
    }

    points := make([]NMAccessPoint, len(accessPaths))

    for i, p := range accessPaths {
        if p != "" {
            points[i] = NMAccessPoint{
                conn: dev.conn,
                Path: p,
            }
        }
    }

    return points
}
