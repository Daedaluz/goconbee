package frame

//go:generate stringer -output=strings.go -type=Command,Status

type Command byte

const (
	CmdDeviceState         = Command(0x07)
	CmdChangeNetworkState  = Command(0x08)
	CmdReadParameter       = Command(0x0A)
	CmdWriteParameter      = Command(0x0B)
	CmdDeviceStateChanged  = Command(0x0E)
	CmdVersion             = Command(0x0D)
	CmdAPSDataRequest      = Command(0x12)
	CmdAPSDataConfirm      = Command(0x04)
	CmdAPSDataIndication   = Command(0x17)
	CmdMacPollIndication   = Command(0x1C)
	CmdMacBeaconIndication = Command(0x1F)
	CmdUpdateBootloader    = Command(0x21)
	CmdUpdateNeighbor      = Command(0x1D)
	CmdGreenPower          = Command(0x19)
)

type Status byte

func (i Status) Error() string {
	return i.String()
}

const (
	StatusSuccess = Status(iota)
	StatusFailure
	StatusBusy
	StatusTimeout
	StatusUnsupported
	StatusError
	StatusNoNetwork
	StatusInvalidValue
)
