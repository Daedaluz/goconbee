package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
	"io"
	"log"
)

// This frame needs further explanation
type GreenPower struct {
	IEEEAddr               uint64
	Seq                    uint16
	FrameType              GPFrameType
	NWKProtocolVersion     uint8
	AutoCommissioning      bool
	NWKExtensionFlag       bool
	ExtApplicationID       uint8
	ExtApplicationSpecific uint8
	GPDSrcID               uint32
	FrameCounter           uint16
	Data                   []byte
}

func (g *GreenPower) CommandID() frame.Command {
	return frame.CmdGreenPower
}

func (g *GreenPower) decode(f frame.Frame) error {
	data := f.Data()
	log.Printf("%X %X %X", data[0:8], data[8:10], data[10:])
	r := bytes.NewReader(data)
	binary.Read(r, binary.LittleEndian, &g.IEEEAddr)
	binary.Read(r, binary.LittleEndian, &g.Seq)
	data, _ = io.ReadAll(r)
	r = bytes.NewReader(data)

	nwkFrameControl, _ := r.ReadByte()
	g.FrameType = GPFrameType(nwkFrameControl & 0b00000011)
	g.NWKProtocolVersion = (nwkFrameControl >> 2) & 0b11
	g.AutoCommissioning = (nwkFrameControl & 0b01000000) > 0
	g.NWKExtensionFlag = (nwkFrameControl & 0b10000000) > 0

	if g.NWKExtensionFlag {
		extFrame, _ := r.ReadByte()
		g.ExtApplicationID = extFrame & 0b00000111
		g.ExtApplicationSpecific = (extFrame & 0b11111000) >> 3
	}

	if g.FrameType == GPFrameTypeData && g.ExtApplicationID == 0 {
		binary.Read(r, binary.LittleEndian, &g.GPDSrcID)
	} else if g.FrameType == GPFrameTypeMaintenance && g.NWKExtensionFlag && g.ExtApplicationID == 0 {
		binary.Read(r, binary.LittleEndian, &g.GPDSrcID)
	}

	if g.NWKExtensionFlag {
		switch g.ExtApplicationID {
		case 0b000, 0b010:
			// GP
			binary.Read(r, binary.LittleEndian, &g.FrameCounter)
		case 0b001:
			// LPED
		}
	}
	g.Data, _ = io.ReadAll(r)
	//	binary.Read(r, binary.LittleEndian, &g.IEEEAddr)
	//	binary.Read(r, binary.LittleEndian, &g.Seq)
	//	binary.Read(r, binary.LittleEndian, &g.Unknown1)
	//	binary.Read(r, binary.LittleEndian, &g.State)
	//	binary.Read(r, binary.LittleEndian, &g.Unknown2)
	//	g.Extra, _ = io.ReadAll(r)
	return nil
}
