package main

import (
	"github.com/daedaluz/goconbee/serial"
	"log"
	"time"
)

var port *serial.Port

func main() {
	var err error

	handler := func(port *serial.Port, msg serial.CommandID) {
		switch x := msg.(type) {
		case *serial.DeviceStateChanged:
			log.Println("DeviceChange:", x)
		case *serial.MacPollIndication:
			log.Println("Poll:", x)
		case *serial.MacBeaconIndication:
			log.Println("Beacon:", x)
		case *serial.GreenPower:
			log.Printf("GreenPower: %v", msg)
		default:
			log.Println("Unknown:", msg)
		}
	}

	port, err = serial.Open("/dev/ttyACM0", &serial.Handlers{
		UnsolicitedHandler: handler,
		DisconnectHandler: func(port2 *serial.Port) {
			log.Println("Disconnected")
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(port.ReadFirmwareVersion())
	PANID := uint16(0)
	ProtocolVersion := uint16(0)
	port.ReadParameter(serial.ParamNWKPANID, &PANID)
	port.ReadParameter(serial.ParamProtocolVersion, &ProtocolVersion)
	NetworkKey, err := port.ReadParameterRaw(serial.ParamNetworkKey)
	Security := serial.SecurityMode(0)
	port.ReadParameter(serial.ParamSecurityMode, &Security)
	log.Printf("%s; Key: %.32X", Security, NetworkKey)

	log.Printf("%.4X %.4X\n", PANID, ProtocolVersion)

	go func() {
		for {
			ticker := time.NewTicker(time.Second * 30)
			port.WriteParameter(serial.ParamWatchdogTTL, uint8(60))
			for _ = range ticker.C {
				port.WriteParameter(serial.ParamWatchdogTTL, uint8(60))
			}
		}
	}()

	go func() {
		for {
			ticker := time.NewTicker(time.Second * 5)
			for _ = range ticker.C {
				state, err := port.GetDeviceState()
				log.Println("State:", state, err)
				if state.DataIndication {
					data, err := port.ReadReceivedData(serial.FlagReadShortSourceAddress)
					log.Println("DATA:", data, err)
				}
				if state.DataConfirm {
					data, err := port.QuerySendData()
					log.Println(data, err)
				}
			}
		}
	}()

	log.Println(port.ReadParameterRaw(serial.ParamNWKPANID))
	state, err := port.GetDeviceState()
	log.Println(state, err)
	if state.DataIndication {
		log.Println(port.ReadReceivedData(serial.FlagReadShortSourceAddress))
	}
	log.Println(port.GetDeviceState())
	log.Println(port.WriteParameter(serial.ParamOpenNetwork, []byte{60}))
	log.Println(port.ReadParameterRaw(serial.ParamOpenNetwork))

	x := &serial.ZDOParameter{}
	port.WriteParameter(serial.ParamZDOSlot, serial.ZDODefaultSlot1, uint8(1))
	port.ReadParameter(serial.ParamZDOSlot, x, uint8(1))
	log.Println(x)

	time.Sleep(time.Minute * 50)
}
