package serial

import (
	"fmt"
	"github.com/daedaluz/goconbee/serial/frame"
	"log"
	"sync"
	"time"
)

type command interface {
	init(c *Port) error
	handle(c *Port, resultCh chan<- any, frame frame.Frame, err error) bool
	ping(now time.Time, resultCh chan<- any)
	finish(c *Port, x any)
	exit(c *Port)
	onError(err error)
}

type commandHandler struct {
	exitCh         chan bool
	doneCh         chan any
	lock           sync.Mutex
	currentCommand command
	c              *Port
}

func (m *commandHandler) ping(now time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.currentCommand != nil {
		m.currentCommand.ping(now, m.doneCh)
	}
}

func newCommandHandler(c *Port) *commandHandler {
	return &commandHandler{
		c:      c,
		lock:   sync.Mutex{},
		exitCh: make(chan bool, 1),
	}
}

func (m *commandHandler) handle(frame frame.Frame, err error) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.currentCommand != nil {
		return m.currentCommand.handle(m.c, m.doneCh, frame, err)
	}
	return false
}

func (m *commandHandler) start() {
	for {
		select {
		case <-m.exitCh:
			return
		case cmd := <-m.c.cmdCh:
			if m.processCommand(cmd) {
				return
			}
		}
	}
}
func (m *commandHandler) processCommand(cmd command) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	m.doneCh = make(chan any)
	defer close(m.doneCh)
	m.lock.Lock()
	m.currentCommand = cmd
	m.lock.Unlock()
	if err := cmd.init(m.c); err != nil {
		cmd.onError(err)
		m.lock.Lock()
		m.currentCommand = nil
		m.lock.Unlock()
		return false
	}
	exit := false
	select {
	case <-m.exitCh:
		m.lock.Lock()
		if m.currentCommand != nil {
			cmd.exit(m.c)
		}
		m.lock.Unlock()
		exit = true
	case x := <-m.doneCh:
		if err, ok := x.(Error); ok {
			cmd.onError(err)
		} else {
			cmd.finish(m.c, x)
		}
	}
	m.lock.Lock()
	m.currentCommand = nil
	m.lock.Unlock()
	return exit
}

type CommandID interface {
	CommandID() frame.Command
}

type request interface {
	CommandID
	encode(seqNumber uint8) frame.Frame
}

type response interface {
	CommandID
	decode(f frame.Frame) error
}

type requestResponseCommand struct {
	start  time.Time
	req    request
	res    response
	seq    uint8
	doneCh chan any
}

func (g *requestResponseCommand) init(c *Port) error {
	g.start = time.Now()
	g.seq = c.getSeqNumber()
	f := g.req.encode(g.seq)
	if err := c.writeFrame(f); err != nil {
		return err
	}
	return nil
}

func (g *requestResponseCommand) handle(c *Port, resultCh chan<- any, f frame.Frame, err error) bool {
	if f.CommandID() == g.req.CommandID() && f.SeqNumber() == g.seq {
		if err := g.res.decode(f); err != nil {
			resultCh <- err
			return true
		}
		resultCh <- g.res
		return true
	}
	return false
}

func (g *requestResponseCommand) ping(now time.Time, resultCh chan<- any) {
	if now.Sub(g.start) > time.Second*7 {
		resultCh <- fmt.Errorf("command timeout")
	}
}

func (g *requestResponseCommand) finish(c *Port, x any) {
	g.doneCh <- x
	close(g.doneCh)
}

func (g *requestResponseCommand) exit(c *Port) {
	g.doneCh <- fmt.Errorf("exit received before finished")
	close(g.doneCh)
}

func (g *requestResponseCommand) onError(err error) {
	g.doneCh <- err
	close(g.doneCh)
}

func (g *requestResponseCommand) wait() (res any) {
	defer func() {
		if err := recover(); err != nil {
			res = nil
		}
	}()
	return <-g.doneCh
}

func newRequestResponseCommand(in request, out response) *requestResponseCommand {
	return &requestResponseCommand{
		req:    in,
		res:    out,
		doneCh: make(chan any),
	}
}
