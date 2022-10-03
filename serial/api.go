package serial

import (
	"bytes"
	"encoding/binary"
)

func (p *Port) ReadFirmwareVersion() (*FirmwareVersion, error) {
	version := &FirmwareVersion{}
	cmd := newRequestResponseCommand(&readFirmwareVersionRequest{}, version)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return version, nil
}

func (p *Port) ReadParameterRaw(param ParameterID) ([]byte, error) {
	paramResp := &ReadParameterResponse{}
	cmd := newRequestResponseCommand(&readParameterRequest{parameterId: param}, paramResp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return paramResp.value, nil
}

func (p *Port) ReadParameter(param ParameterID, out any, args ...any) error {
	paramResp := &ReadParameterResponse{}

	buff := &bytes.Buffer{}
	for _, arg := range args {
		binary.Write(buff, binary.LittleEndian, arg)
	}

	cmd := newRequestResponseCommand(&readParameterRequest{parameterId: param, args: buff.Bytes()}, paramResp)
	p.cmdCh <- cmd
	res := cmd.wait()

	if o, ok := out.(paramDecoder); ok {
		o.decode(paramResp.value)
	} else {
		r := bytes.NewReader(paramResp.value)
		binary.Read(r, binary.LittleEndian, out)
	}
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

func (p *Port) WriteParameterRaw(param ParameterID, value []byte) error {
	resp := &WriteParameterResponse{}
	cmd := newRequestResponseCommand(&writeParameterRequest{
		parameterID: param,
		value:       value,
	}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

func (p *Port) WriteParameter(param ParameterID, value any, args ...any) error {
	buff := &bytes.Buffer{}
	for _, arg := range args {
		binary.Write(buff, binary.LittleEndian, arg)
	}
	if encoder, ok := value.(paramEncoder); ok {
		buff.Write(encoder.encode())
	} else {
		binary.Write(buff, binary.LittleEndian, value)
	}
	resp := &WriteParameterResponse{}
	cmd := newRequestResponseCommand(&writeParameterRequest{
		parameterID: param,
		value:       buff.Bytes(),
	}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

func (p *Port) GetDeviceState() (*DeviceState, error) {
	stateResp := &DeviceState{}
	cmd := newRequestResponseCommand(&deviceStateRequest{}, stateResp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return stateResp, nil
}

func (p *Port) ChangeNetworkState(state NetworkState) error {
	resp := &ChangeNetworkStateResponse{}
	cmd := newRequestResponseCommand(&changeNetworkStateRequest{NetworkState: state}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

func (p *Port) ReadReceivedData(flags ReadDataFlag) (*ApsData, error) {
	resp := &ApsData{}
	cmd := newRequestResponseCommand(&apsReadDataRequest{flags: flags}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return resp, nil
}

func (p *Port) SendData(reqID uint8, dstAddr Address, profileID, clusterID uint16, srcEP uint8, data []byte, opts TXOptions, radius uint8, srcRoute ...uint16) (*SendDataResponse, error) {
	resp := &SendDataResponse{}
	req := &enqueueSendDataRequest{
		RequestID:  reqID,
		DstAddress: dstAddr,
		ProfileID:  profileID,
		ClusterID:  clusterID,
		SrcEP:      srcEP,
		Data:       data,
		Options:    opts,
		Radius:     radius,
	}
	if len(srcRoute) > 0 {
		req.Flags |= sendDataFlagSourceRouting
		req.Relay = srcRoute
	}
	cmd := newRequestResponseCommand(req, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return resp, err
	}
	return resp, nil
}

func (p *Port) QuerySendData() (*QuerySendDataResponse, error) {
	resp := &QuerySendDataResponse{}
	cmd := newRequestResponseCommand(&querySendDataRequest{}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return resp, nil
}

// Currently not working.. why?
func (p *Port) AddNeighbor(nwk uint16, IEEEAddr uint64, macCapabilities MacCapabilities) error {
	resp := &UpdateNeighborResponse{}
	cmd := newRequestResponseCommand(&updateNeighborRequest{
		Action:          actionAdd,
		NWK:             nwk,
		IEEEAddr:        IEEEAddr,
		MacCapabilities: macCapabilities,
	}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

// Currently not working.. why?
func (p *Port) RemoveNeighbor(nwk uint16, IEEEAddr uint64, macCapabilities MacCapabilities) error {
	resp := &UpdateNeighborResponse{}
	cmd := newRequestResponseCommand(&updateNeighborRequest{
		Action:          actionRemove,
		NWK:             nwk,
		IEEEAddr:        IEEEAddr,
		MacCapabilities: macCapabilities,
	}, resp)
	p.cmdCh <- cmd
	res := cmd.wait()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}
