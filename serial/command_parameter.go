package serial

import (
	"bytes"
	"encoding/binary"
	"github.com/daedaluz/goconbee/serial/frame"
)

type readParameterRequest struct {
	parameterId ParameterID
	args        []byte
}

func (r *readParameterRequest) CommandID() frame.Command {
	return frame.CmdReadParameter
}

func (r *readParameterRequest) encode(seqNumber uint8) frame.Frame {
	payload := []byte{0, 0, uint8(r.parameterId)}
	payload = append(payload, r.args...)
	binary.LittleEndian.PutUint16(payload[0:2], 1+uint16(len(r.args)))
	return frame.NewFrame(frame.CmdReadParameter, seqNumber, payload)
}

type ReadParameterResponse struct {
	parameterID ParameterID
	value       []byte
}

func (r *ReadParameterResponse) CommandID() frame.Command {
	return frame.CmdReadParameter
}

func (r *ReadParameterResponse) decode(f frame.Frame) error {
	data := f.Data()
	r.parameterID = ParameterID(data[0])
	r.value = data[3:]
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}

type writeParameterRequest struct {
	parameterID ParameterID
	value       []byte
}

func (w *writeParameterRequest) CommandID() frame.Command {
	return frame.CmdWriteParameter
}

func (w *writeParameterRequest) encode(seqNumber uint8) frame.Frame {
	buff := &bytes.Buffer{}
	binary.Write(buff, binary.LittleEndian, uint16(1+len(w.value)))
	buff.WriteByte(byte(w.parameterID))
	buff.Write(w.value)
	return frame.NewFrame(frame.CmdWriteParameter, seqNumber, buff.Bytes())
}

type WriteParameterResponse struct {
	parameterID ParameterID
}

func (w *WriteParameterResponse) CommandID() frame.Command {
	return frame.CmdWriteParameter
}

func (w *WriteParameterResponse) decode(f frame.Frame) error {
	status := f.Status()
	w.parameterID = ParameterID(f.Data()[1])
	if status != frame.StatusSuccess {
		return status
	}
	return nil
}
