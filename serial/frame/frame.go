package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func crc16(data []byte) uint16 {
	crc := uint16(0)
	for _, x := range data {
		crc += uint16(x) & 0x00FF
	}
	crc = (crc ^ 0xFFFF) + 1
	return crc
}

type Frame []byte

func (f Frame) getCRC() uint16 {
	crc := f[len(f)-2:]
	return binary.LittleEndian.Uint16(crc)
}

func (f Frame) CheckCRC() bool {
	if len(f) < 3 {
		return false
	}
	crcData := f[0 : len(f)-2]
	calc := crc16(crcData)
	return calc == f.getCRC()
}

func (f Frame) CommandID() Command {
	return Command(f[0])
}

func (f Frame) SeqNumber() uint8 {
	return f[1]
}

func (f Frame) Status() Status {
	return Status(f[2])
}

func (f Frame) Data() []byte {
	idxStart := 5
	idxEnd := len(f) - 2
	return f[idxStart:idxEnd]
}

func NewFrame(cmd Command, seq uint8, payload []byte) Frame {
	x := bytes.NewBuffer(make([]byte, 0, 5+len(payload)+2))
	x.WriteByte(byte(cmd)) // Command
	x.WriteByte(seq)       // sequence number
	x.WriteByte(0)         // reserved

	_ = binary.Write(x, binary.LittleEndian, uint16(len(payload)+5)) // +5 for the bytes we already wrote, including payload length
	x.Write(payload)
	crc := crc16(x.Bytes())
	_ = binary.Write(x, binary.LittleEndian, crc)
	return x.Bytes()
}

func (f Frame) String() string {
	return fmt.Sprintf("[%s seq:%d s:%s %X]", f.CommandID(), f.SeqNumber(), f.Status(), f.Data())
}
