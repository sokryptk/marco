package repository

type DeviceType int
type ConnectionStatus uint32

const (
    DeviceTypeUnknown DeviceType = 0
    DeviceTypeEthernet DeviceType = 1
    DeviceTypeWifi DeviceType = 2
    DeviceTypeBridge DeviceType = 13
    DeviceTypeTunnel DeviceType = 16
    DeviceTypeBluetooth DeviceType = 5
)

const (
    ConnectionStatusNone ConnectionStatus = 0
    ConnectionStatusActivated ConnectionStatus = 100
    ConnectionStatusErrNoSecrets ConnectionStatus = 7
    ConnectionStatusErrFailed ConnectionStatus = 120
)

func (status ConnectionStatus) Equal(v uint32) bool {
    return uint32(status) == v
}

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
    Connect(options ConnectOptions) ConnectionStatus
}

type ConnectOptions struct {
    Password *string
}

