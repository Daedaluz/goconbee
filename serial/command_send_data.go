package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
)

type TXOptions uint8

const TXOptUseAPSAck = TXOptions(0x04)

type enqueueSendDataRequest struct {
	Flags      sendDataFlags
	RequestID  uint8
	DstAddress Address
	ProfileID  uint16
	ClusterID  uint16
	SrcEP      uint8
	Data       []byte
	Options    TXOptions
	Radius     uint8
	Relay      []uint16
}

func (e *enqueueSendDataRequest) CommandID() frame.Command {
	return frame.CmdAPSDataRequest
}

func (e *enqueueSendDataRequest) encode(seqNumber uint8) frame.Frame {
	payloadBuff := &bytes.Buffer{}
	payloadBuff.WriteByte(e.RequestID)
	payloadBuff.WriteByte(byte(e.Flags))
	payloadBuff.WriteByte(byte(e.DstAddress.Mode))
	switch e.DstAddress.Mode {
	case AddressGroup, AddressNWK:
		binary.Write(payloadBuff, binary.LittleEndian, e.DstAddress.Short)
	case AddressIEEE:
		binary.Write(payloadBuff, binary.LittleEndian, e.DstAddress.Extended)
	case AddressNWKAndIEEE:
		binary.Write(payloadBuff, binary.LittleEndian, e.DstAddress.Short)
		binary.Write(payloadBuff, binary.LittleEndian, e.DstAddress.Extended)
	}
	switch e.DstAddress.Mode {
	case AddressNWK, AddressIEEE, AddressNWKAndIEEE:
		payloadBuff.WriteByte(e.DstAddress.Endpoint)
	}
	binary.Write(payloadBuff, binary.LittleEndian, e.ProfileID)
	binary.Write(payloadBuff, binary.LittleEndian, e.ClusterID)
	binary.Write(payloadBuff, binary.LittleEndian, e.SrcEP)
	binary.Write(payloadBuff, binary.LittleEndian, uint16(len(e.Data)))
	payloadBuff.Write(e.Data)
	payloadBuff.WriteByte(byte(e.Options))
	payloadBuff.WriteByte(e.Radius)

	if e.Flags&0x02 > 0 {
		payloadBuff.WriteByte(byte(len(e.Relay)))
		for _, x := range e.Relay {
			binary.Write(payloadBuff, binary.LittleEndian, x)
		}
	}

	dataBuff := &bytes.Buffer{}
	payload := payloadBuff.Bytes()
	binary.Write(dataBuff, binary.LittleEndian, uint16(len(payload)))
	dataBuff.Write(payload)
	return frame.NewFrame(frame.CmdAPSDataRequest, seqNumber, dataBuff.Bytes())
}

type SendDataResponse struct {
	NetworkState         NetworkState
	DataConfirm          bool
	DataIndication       bool
	ConfigurationChanged bool
	FreeSlots            bool
	RequestID            uint8
}

func (e *SendDataResponse) CommandID() frame.Command {
	return frame.CmdAPSDataRequest
}

func (e *SendDataResponse) decode(f frame.Frame) error {
	data := f.Data()
	b := data[2]
	e.NetworkState = NetworkState(b & 0b00000011)
	e.DataConfirm = b&0b00000100 > 0
	e.DataIndication = b&0b00001000 > 0
	e.ConfigurationChanged = b&0b00010000 > 0
	e.FreeSlots = b&0b00100000 > 0
	e.RequestID = data[3]
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
