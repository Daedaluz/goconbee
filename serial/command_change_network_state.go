package serial

import (
	"github.com/daedaluz/goconbee/serial/frame"
)

type changeNetworkStateRequest struct {
	NetworkState NetworkState
}

func (c *changeNetworkStateRequest) CommandID() frame.Command {
	return frame.CmdChangeNetworkState
}

func (c *changeNetworkStateRequest) encode(seqNumber uint8) frame.Frame {
	return frame.NewFrame(frame.CmdChangeNetworkState, seqNumber, []byte{byte(c.NetworkState)})
}

type ChangeNetworkStateResponse struct {
	NetworkState NetworkState
}

func (c *ChangeNetworkStateResponse) CommandID() frame.Command {
	return frame.CmdChangeNetworkState
}

func (c *ChangeNetworkStateResponse) decode(f frame.Frame) error {
	c.NetworkState = NetworkState(f.Data()[0])
	if f.Status() != frame.StatusSuccess {
		return f.Status()
	}
	return nil
}
