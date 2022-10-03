package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
)

type updateNeighborRequest struct {
	Action          updateNeighborAction
	NWK             uint16
	IEEEAddr        uint64
	MacCapabilities MacCapabilities
}

func (a *updateNeighborRequest) CommandID() frame.Command {
	return frame.CmdUpdateNeighbor
}

func (a *updateNeighborRequest) encode(seqNumber uint8) frame.Frame {
	buff := &bytes.Buffer{}
	buff.WriteByte(byte(a.Action))
	binary.Write(buff, binary.LittleEndian, a.NWK)
	binary.Write(buff, binary.LittleEndian, a.IEEEAddr)
	buff.WriteByte(byte(a.MacCapabilities))
	return frame.NewFrame(frame.CmdUpdateNeighbor, seqNumber, buff.Bytes())
}

type UpdateNeighborResponse struct {
	Data []byte
}

func (a *UpdateNeighborResponse) CommandID() frame.Command {
	return frame.CmdUpdateNeighbor
}

func (a *UpdateNeighborResponse) decode(f frame.Frame) error {
	a.Data = f.Data()
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
