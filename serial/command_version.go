package serial

import (
	"github.com/daedaluz/goconbee/serial/frame"
)

type readFirmwareVersionRequest struct {
}

func (r *readFirmwareVersionRequest) CommandID() frame.Command {
	return frame.CmdVersion
}

func (r *readFirmwareVersionRequest) encode(seqNumber uint8) frame.Frame {
	return frame.NewFrame(frame.CmdVersion, seqNumber, []byte{0, 0, 0, 0})
}

type FirmwareVersion struct {
	Major    uint8
	Minor    uint8
	Platform Platform
}

func (r *FirmwareVersion) CommandID() frame.Command {
	return frame.CmdVersion
}

func (r *FirmwareVersion) decode(f frame.Frame) error {
	data := f.Data()
	r.Major = data[3]
	r.Minor = data[2]
	r.Platform = Platform(data[1])
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
