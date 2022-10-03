package serial

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type paramEncoder interface {
	encode() []byte
}

type paramDecoder interface {
	decode(data []byte)
}

type ZDOParameter struct {
	Endpoint      uint8
	ProfileID     uint16
	DeviceID      uint16
	DeviceVersion uint8

	InClusters  []uint16
	OutClusters []uint16
}

var (
	ZDODefaultSlot0 = &ZDOParameter{
		Endpoint:      0x01,
		ProfileID:     0x0104,
		DeviceID:      0x0005,
		DeviceVersion: 0x01,
		InClusters:    []uint16{0x0000, 0x0006, 0x000A, 0x0019, 0x0501},
		OutClusters:   []uint16{0x0001, 0x0020, 0x0500, 0x0502},
	}
	ZDODefaultSlot1 = &ZDOParameter{
		Endpoint:      0xF2,
		ProfileID:     0xA1E0,
		DeviceID:      0x0064,
		DeviceVersion: 0x01,
		InClusters:    []uint16{},
		OutClusters:   []uint16{0x0021},
	}
)

func (z *ZDOParameter) String() string {
	return fmt.Sprintf("{ep:0x%.2x profile:0x%.4x deviceid:0x%.4x devicever:0x%.4x servers:%.4x clients:%.4x}",
		z.Endpoint, z.ProfileID, z.DeviceID, z.DeviceVersion, z.InClusters, z.OutClusters)
}

func (z *ZDOParameter) encode() []byte {
	buff := &bytes.Buffer{}
	buff.WriteByte(z.Endpoint)
	binary.Write(buff, binary.LittleEndian, z.ProfileID)
	binary.Write(buff, binary.LittleEndian, z.DeviceID)
	buff.WriteByte(z.DeviceVersion)
	buff.WriteByte(byte(len(z.InClusters)))
	for _, cluster := range z.InClusters {
		binary.Write(buff, binary.LittleEndian, cluster)
	}
	buff.WriteByte(byte(len(z.OutClusters)))
	for _, cluster := range z.OutClusters {
		binary.Write(buff, binary.LittleEndian, cluster)
	}
	return buff.Bytes()
}

func (z *ZDOParameter) decode(data []byte) {
	r := bytes.NewReader(data)
	_, _ = r.ReadByte()
	z.Endpoint, _ = r.ReadByte()
	binary.Read(r, binary.LittleEndian, &z.ProfileID)
	binary.Read(r, binary.LittleEndian, &z.DeviceID)
	z.DeviceVersion, _ = r.ReadByte()
	nIn, _ := r.ReadByte()
	for i := uint8(0); i < nIn; i++ {
		cluster := uint16(0)
		binary.Read(r, binary.LittleEndian, &cluster)
		z.InClusters = append(z.InClusters, cluster)
	}
	nOut, _ := r.ReadByte()
	for i := uint8(0); i < nOut; i++ {
		cluster := uint16(0)
		binary.Read(r, binary.LittleEndian, &cluster)
		z.OutClusters = append(z.OutClusters, cluster)
	}
}
