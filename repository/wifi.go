package repository

type DeviceType int

const (
    DeviceTypeUnknown DeviceType = 0
    DeviceTypeEthernet DeviceType = 1
    DeviceTypeWifi DeviceType = 2
    DeviceTypeBluetooth DeviceType = 5
)

type WiFi interface {
    GetDevices() []Device
}

type Device interface {
    GetHwAddress() string
    GetAccessPoints() []AccessPoint
    GetDeviceType() DeviceType
    RequestScan() error
}

type AccessPoint interface {
    GetSSID() string
    GetFrequency() uint
    GetStrength() uint

    Connect() error
}