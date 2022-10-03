package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
	"io"
)

type querySendDataRequest struct {
}

func (q *querySendDataRequest) CommandID() frame.Command {
	return frame.CmdAPSDataConfirm
}

func (q *querySendDataRequest) encode(seqNumber uint8) frame.Frame {
	return frame.NewFrame(frame.CmdAPSDataConfirm, seqNumber, nil)
}

type QuerySendDataResponse struct {
	NetworkState         NetworkState
	DataConfirm          bool
	DataIndication       bool
	ConfigurationChanged bool
	FreeSlots            bool

	RequestID  uint8
	DstAddress Address
	SrcEP      uint8
	Status     uint8
}

func (q *QuerySendDataResponse) CommandID() frame.Command {
	return frame.CmdAPSDataConfirm
}

func (q *QuerySendDataResponse) decode(f frame.Frame) error {
	r := bytes.NewReader(f.Data())
	r.Seek(2, io.SeekCurrent)
	b, _ := r.ReadByte()
	q.NetworkState = NetworkState(b & 0b00000011)
	q.DataConfirm = b&0b00000100 > 0
	q.DataIndication = b&0b00001000 > 0
	q.ConfigurationChanged = b&0b00010000 > 0
	q.FreeSlots = b&0b00100000 > 0

	binary.Read(r, binary.LittleEndian, &q.DstAddress.Mode)
	switch q.DstAddress.Mode {
	case AddressGroup, AddressNWK:
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Short)
		if q.DstAddress.Mode == AddressNWK {
			binary.Read(r, binary.LittleEndian, &q.DstAddress.Endpoint)
		}
	case AddressIEEE:
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Extended)
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Endpoint)
	case AddressNWKAndIEEE:
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Short)
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Extended)
		binary.Read(r, binary.LittleEndian, &q.DstAddress.Endpoint)
	}
	binary.Read(r, binary.LittleEndian, &q.SrcEP)
	binary.Read(r, binary.LittleEndian, &q.Status)
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
