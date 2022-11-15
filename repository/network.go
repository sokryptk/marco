package repository

type DeviceType int

const (
    DeviceTypeUnknown DeviceType = 0
    DeviceTypeEthernet DeviceType = 1
    DeviceTypeWifi DeviceType = 2
    DeviceTypeBridge DeviceType = 13
    DeviceTypeTunnel DeviceType = 16
    DeviceTypeBluetooth DeviceType = 5
)

type Network interface {
    GetDevices() []Device
    Close() error
}

type Device interface {
    GetHwAddress() string
GetDeviceType() DeviceType
}

type WiFiDevice interface {
    Device
    GetHwAddress() string
    GetAccessPoints() []AccessPoint
    GetActiveConnection() (AccessPoint, error)
    RequestScan() error
}

type EthernetDevice interface {
    Device
}

type AccessPoint interface {
    GetSSID() string
    GetFrequency() uint
    GetStrength() uint
    Connect() error
}
