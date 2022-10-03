package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
	"io"
)

type MacPollIndication struct {
	SrcAddr Address
	LQI     uint8
	RSSI    int8
	Extra   []byte
}

func (m *MacPollIndication) CommandID() frame.Command {
	return frame.CmdMacPollIndication
}

func (m *MacPollIndication) decode(f frame.Frame) error {
	r := bytes.NewReader(f.Data())
	r.Seek(2, io.SeekCurrent)
	binary.Read(r, binary.LittleEndian, &m.SrcAddr.Mode)
	switch m.SrcAddr.Mode {
	case AddressNWK:
		binary.Read(r, binary.LittleEndian, &m.SrcAddr.Short)
	case AddressIEEE:
		binary.Read(r, binary.LittleEndian, &m.SrcAddr.Extended)
	}
	binary.Read(r, binary.LittleEndian, &m.LQI)
	binary.Read(r, binary.LittleEndian, &m.RSSI)
	m.Extra, _ = io.ReadAll(r)
	return nil
}

type MacBeaconIndication struct {
	SrcAddr  uint16
	PANID    uint16
	Channel  uint8
	Flags    uint8
	UpdateID uint8
	Extra    []byte
}

func (m *MacBeaconIndication) CommandID() frame.Command {
	return frame.CmdMacBeaconIndication
}

func (m *MacBeaconIndication) decode(f frame.Frame) error {
	r := bytes.NewReader(f.Data())
	binary.Read(r, binary.LittleEndian, &m.SrcAddr)
	binary.Read(r, binary.LittleEndian, &m.PANID)
	binary.Read(r, binary.LittleEndian, &m.Channel)
	binary.Read(r, binary.LittleEndian, &m.Flags)
	binary.Read(r, binary.LittleEndian, &m.UpdateID)
	m.Extra, _ = io.ReadAll(r)
	return nil
}
