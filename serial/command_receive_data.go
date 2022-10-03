package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
	"io"
)

type apsReadDataRequest struct {
	flags ReadDataFlag
}

func (a *apsReadDataRequest) CommandID() frame.Command {
	return frame.CmdAPSDataIndication
}

func (a *apsReadDataRequest) encode(seqNumber uint8) frame.Frame {
	var params = []byte{0, 0}

	if a.flags != 0 {
		binary.LittleEndian.PutUint16(params[:], 1)
		params = append(params, byte(a.flags))
	}
	f := frame.NewFrame(frame.CmdAPSDataIndication, seqNumber, params)
	return f
}

type ApsData struct {
	NetworkState         NetworkState
	DataConfirm          bool
	DataIndication       bool
	ConfigurationChanged bool
	FreeSlots            bool

	DstAddress Address
	SrcAddress Address

	ProfileID uint16
	ClusterID uint16
	Data      []byte

	LastHop uint16
	LQI     uint8
	RSSI    int8
}

func (a *ApsData) CommandID() frame.Command {
	return frame.CmdAPSDataIndication
}

func (a *ApsData) decode(f frame.Frame) error {
	data := f.Data()
	r := bytes.NewReader(data)
	r.Seek(2, io.SeekCurrent)
	b, _ := r.ReadByte()
	a.NetworkState = NetworkState(b & 0b00000011)
	a.DataConfirm = b&0b00000100 > 0
	a.DataIndication = b&0b00001000 > 0
	a.ConfigurationChanged = b&0b00010000 > 0
	a.FreeSlots = b&0b00100000 > 0

	binary.Read(r, binary.LittleEndian, &a.DstAddress.Mode)
	switch a.DstAddress.Mode {
	case AddressGroup, AddressNWK:
		binary.Read(r, binary.LittleEndian, &a.DstAddress.Short)
	case AddressIEEE:
		binary.Read(r, binary.LittleEndian, &a.DstAddress.Extended)
	case AddressNWKAndIEEE:
		binary.Read(r, binary.LittleEndian, &a.DstAddress.Short)
		binary.Read(r, binary.LittleEndian, &a.DstAddress.Extended)
	}
	binary.Read(r, binary.LittleEndian, &a.DstAddress.Endpoint)

	binary.Read(r, binary.LittleEndian, &a.SrcAddress.Mode)
	switch a.DstAddress.Mode {
	case AddressGroup, AddressNWK:
		binary.Read(r, binary.LittleEndian, &a.SrcAddress.Short)
	case AddressIEEE:
		binary.Read(r, binary.LittleEndian, &a.SrcAddress.Extended)
	case AddressNWKAndIEEE:
		binary.Read(r, binary.LittleEndian, &a.SrcAddress.Short)
		binary.Read(r, binary.LittleEndian, &a.SrcAddress.Extended)
	}
	binary.Read(r, binary.LittleEndian, &a.SrcAddress.Endpoint)

	binary.Read(r, binary.LittleEndian, &a.ProfileID)
	binary.Read(r, binary.LittleEndian, &a.ClusterID)
	asduLen := uint16(0)
	binary.Read(r, binary.LittleEndian, &asduLen)

	a.Data = make([]byte, asduLen)
	r.Read(a.Data)

	binary.Read(r, binary.LittleEndian, &a.LastHop)
	a.LQI, _ = r.ReadByte()
	r.Seek(4, io.SeekCurrent)
	binary.Read(r, binary.LittleEndian, &a.RSSI)
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
