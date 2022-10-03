package serial

import (
	"github.com/daedaluz/goconbee/serial/frame"
)

type deviceStateRequest struct {
}

func (d *deviceStateRequest) CommandID() frame.Command {
	return frame.CmdDeviceState
}

func (d *deviceStateRequest) encode(seqNumber uint8) frame.Frame {
	return frame.NewFrame(frame.CmdDeviceState, seqNumber, []byte{0, 0, 0})
}

type DeviceState struct {
	NetworkState         NetworkState
	DataConfirm          bool
	DataIndication       bool
	ConfigurationChanged bool
	FreeSlots            bool
}

func (d *DeviceState) CommandID() frame.Command {
	return frame.CmdDeviceState
}

func (d *DeviceState) decode(f frame.Frame) error {
	data := f.Data()
	d.NetworkState = NetworkState(data[0] & 0b00000011)
	d.DataConfirm = data[0]&0b00000100 > 0
	d.DataIndication = data[0]&0b00001000 > 0
	d.ConfigurationChanged = data[0]&0b00010000 > 0
	d.FreeSlots = data[0]&0b00100000 > 0
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}

type DeviceStateChanged struct {
	NetworkState         NetworkState
	DataConfirm          bool
	DataIndication       bool
	ConfigurationChanged bool
	FreeSlots            bool
}

func (d *DeviceStateChanged) CommandID() frame.Command {
	return frame.CmdDeviceStateChanged
}

func (d *DeviceStateChanged) decode(f frame.Frame) error {
	data := f.Data()
	d.NetworkState = NetworkState(data[0] & 0b00000011)
	d.DataConfirm = data[0]&0b00000100 > 0
	d.DataIndication = data[0]&0b00001000 > 0
	d.ConfigurationChanged = data[0]&0b00010000 > 0
	d.FreeSlots = data[0]&0b00100000 > 0
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
