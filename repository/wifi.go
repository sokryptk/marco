package repository

type WiFi interface {
    GetDevices() []Device
}

type Device interface {
    GetAccessPoints() []AccessPoint
    RequestScan() error
}

type AccessPoint interface {
    GetSSID() string
    GetFrequency() uint
    GetStrength() uint

    Connect() error
}