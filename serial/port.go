package serial

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/daedaluz/fdev/poll"
	"github.com/daedaluz/goconbee/serial/frame"
	"github.com/daedaluz/goserial"
	"github.com/daedaluz/goslip"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

type UnsolicitedHandler func(port *Port, msg CommandID)

type Address struct {
	Mode     AddressMode
	Endpoint uint8
	Short    uint16
	Extended uint64
}

func (a Address) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("(%d)%s", a.Endpoint, a.Mode))
	switch a.Mode {
	case AddressGroup, AddressNWK:
		parts = append(parts, fmt.Sprintf("0x%.4x", a.Short))
	case AddressIEEE:
		parts = append(parts, fmt.Sprintf("0x%.16x", a.Extended))
	case AddressNWKAndIEEE:
		parts = append(parts, fmt.Sprintf("0x%.16x(%.4x)", a.Extended, a.Short))
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, " "))
}

type Port struct {
	rs232       *serial.Port
	buf         *bufio.Reader
	rw          slip.ReadWriter
	cmdCh       chan command
	handlers    *Handlers
	cmdHandlers []*commandHandler
	seq         atomic.Uint64
	lastPoll    time.Time
}

type DisconnectHandler func(p *Port)

type Handlers struct {
	UnsolicitedHandler
	DisconnectHandler
}

func Open(path string, handlers *Handlers) (*Port, error) {
	options := serial.NewOptions()
	options.SetReadTimeout(time.Second)
	p, err := serial.Open(path, options)
	if err != nil {
		return nil, err
	}
	attrs, err := p.GetAttr()
	if err != nil {
		p.Close()
		return nil, err
	}
	attrs.MakeRaw()
	attrs.SetSpeed(serial.B115200)
	if err := p.SetAttr(serial.TCSANOW, attrs); err != nil {
		p.Close()
		return nil, err
	}

	br := bufio.NewReader(p)
	rw := slip.NewReadWriter2(br, p)

	defaultUnsolicitedHandler := func(p *Port, f CommandID) {
		log.Println("Unhandled: ", f)
	}
	defaultDisconnectHandler := func(p *Port) {
		log.Println("Conbee disconnected")
	}
	if handlers == nil {
		handlers = &Handlers{
			UnsolicitedHandler: defaultUnsolicitedHandler,
			DisconnectHandler:  defaultDisconnectHandler,
		}
	}
	if handlers.UnsolicitedHandler == nil {
		handlers.UnsolicitedHandler = defaultUnsolicitedHandler
	}
	if handlers.DisconnectHandler == nil {
		handlers.DisconnectHandler = defaultDisconnectHandler
	}

	port := &Port{
		rs232:       p,
		buf:         br,
		rw:          rw,
		handlers:    handlers,
		cmdCh:       make(chan command, 100),
		cmdHandlers: make([]*commandHandler, 0, 10),
	}
	port.cmdHandlers = append(port.cmdHandlers, newCommandHandler(port))
	for _, handler := range port.cmdHandlers {
		go handler.start()
	}
	go port.rx()
	return port, nil
}

func (p *Port) SetNHandlers(n int) {
	if len(p.cmdHandlers) == n {
		return
	}

	newHandlers := make([]*commandHandler, n)

	if n < len(p.cmdHandlers) {
		for i, handler := range p.cmdHandlers {
			if i < n {
				newHandlers[i] = p.cmdHandlers[i]
			} else {
				handler.exitCh <- true
			}
		}
		return
	}
	for x := len(p.cmdHandlers); x < n; x++ {
		handler := newCommandHandler(p)
		go handler.start()
		p.cmdHandlers = append(p.cmdHandlers, handler)
	}
}

func (p *Port) readFrame() (frame.Frame, error) {
	for {
		x, err := p.rw.ReadPacket()
		if err != nil {
			return nil, err
		}
		if len(x) != 0 {
			return x, nil
		}
	}
}

func (p *Port) writeFrame(f frame.Frame) error {
	p.rs232.Write([]byte{0300})
	return p.rw.WritePacket(f)
}

func (p *Port) rx() {
	var f frame.Frame
	var err error
outerLoop:
	for f, err = p.readFrame(); err == nil || errors.Is(err, poll.ErrTimeout); f, err = p.readFrame() {
		now := time.Now()
		if now.Sub(p.lastPoll) > time.Second {
			for _, handler := range p.cmdHandlers {
				handler.ping(now)
			}
		}
		if errors.Is(err, poll.ErrTimeout) {
			continue
		}
		if !f.CheckCRC() {
			continue
		}
		for _, handler := range p.cmdHandlers {
			if handler.handle(f, nil) {
				continue outerLoop
			}
		}
		x := CommandID(f)
		switch f.CommandID() {
		case frame.CmdMacBeaconIndication:
			tmp := &MacBeaconIndication{}
			tmp.decode(f)
			x = CommandID(tmp)
		case frame.CmdMacPollIndication:
			tmp := &MacPollIndication{}
			tmp.decode(f)
			x = CommandID(tmp)
		case frame.CmdDeviceStateChanged:
			tmp := &DeviceStateChanged{}
			tmp.decode(f)
			x = CommandID(tmp)
		case frame.CmdGreenPower:
			tmp := &GreenPower{}
			tmp.decode(f)
			x = CommandID(tmp)
		}
		if p.handlers.UnsolicitedHandler != nil {
			p.handlers.UnsolicitedHandler(p, x)
		}
	}
	p.handlers.DisconnectHandler(p)
	p.Close()
}

func (p *Port) getSeqNumber() uint8 {
	x := p.seq.Add(1) % 255
	return uint8(x)
}

func (p *Port) Close() error {
	for _, handler := range p.cmdHandlers {
		select {
		case handler.exitCh <- true:
		default:
		}
	}
	return p.rs232.Close()
}
