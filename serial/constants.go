package serial

//go:generate stringer -output=strings.go -type=Platform,NetworkState,ParameterID,AddressMode,SecurityMode,PANIDMode,MacCapabilities

type Platform byte

const (
	Conbee  = Platform(0x05)
	Conbee2 = Platform(0x07)
)

type NetworkState byte

const (
	NetOffline = NetworkState(iota)
	NetJoining
	NetConnected
	NetLeaving
)

type ParameterID byte

const (
	// R  | u64
	ParamMACAddress = ParameterID(0x01)

	// RW | u16
	ParamNWKPANID = ParameterID(0x05)

	// R  | u16
	ParamNWKAddress = ParameterID(0x07)

	// R  | u64
	ParamNWKExtendedPANID = ParameterID(0x08)

	// RW | u8 0x00 - Router, the node will join a network.
	//         0x01 - Coordinator, the node will form a network and let others nodes join.
	ParamAPSDesignedCoordinator = ParameterID(0x09)

	// RW | u32
	ParamChannelMask = ParameterID(0x0A)

	// RW | u64
	ParamAPSExtendedPANID = ParameterID(0x0B)

	// RW | u64
	ParamTrustCenterAddress = ParameterID(0x0E)

	// RW | u8 0x00 - No security
	//         0x01 - Preconfigured network key
	//         0x02 - Network key from trust center
	//         0x03 - No master but trust center link key
	ParamSecurityMode = ParameterID(0x10)

	// RW | u8:slot u8:endpoint u16:profileId u16:deviceId u8:deviceVersion u8:nIn []u16:InClusters u8:nOut []u16:OutClusters
	ParamZDOSlot = ParameterID(0x13)

	// RW | u8 0x00 - not predefined: The NWK PAN ID will be selected or obtained dynamically.
	//         0x01 - predefined: The value of parameter NWK PANID (0x05) will be used to join or form a network
	ParamPredefinedNWKPANID = ParameterID(0x15)

	// RW | u8[16] Encryption key to secure network traffic
	ParamNetworkKey = ParameterID(0x18)

	// RW | u64    MAC address
	//      u8[16] Link key to be used during joining. The key is only included in a write request and read response. The read request shall only contain the MAC Address.
	ParamLinkKey = ParameterID(0x19)

	// R  | u8   11-26
	ParamCurrentChannel = ParameterID(0x1C)

	// RW | u8 0x00-0xFF
	//         0x00 - Closed
	//         0xFF - Open
	//         Other: Number of seconds to permit joins
	ParamOpenNetwork = ParameterID(0x21)

	// R  | u16  Version of the implemented protocol.
	ParamProtocolVersion = ParameterID(0x22)

	// RW | u8   0-255
	ParamNWKUpdateID = ParameterID(0x24)

	// RW | u32  Watchdog timeout in seconds. Must be reset by the application periodically (since protocol version 0x0108). By writing a lower value like 2 seconds, the firmware can be rebooted.
	ParamWatchdogTTL = ParameterID(0x26)

	// RW | u32  Outgoing security frame messageID. It shall be only set initially when joining or forming a network.
	ParamNWKFrameCounter = ParameterID(0x27)

	// RW | u16 A bitmap describing which ZDP responses the application wants to handle.
	//          The bitmap is not persistant and resets on every power-up of the firmware.
	//          Default value is 0x0000.
	// Supported Flags:
	//   * 0x0001 - Node Descriptor response.
	ParamAppZDPHandling = ParameterID(0x28)
)

type AppZDPHandlingFlag uint16

const (
	AppZDPHandleNodeDescriptor = AppZDPHandlingFlag(0x0001)
)

type AddressMode byte

const (
	AddressGroup = AddressMode(iota + 1)
	AddressNWK
	AddressIEEE
	AddressNWKAndIEEE
)

type ReadDataFlag byte

const (
	FlagReadShortSourceAddress         = ReadDataFlag(0x01)
	FlagLastHop                        = ReadDataFlag(0x02)
	FlagIncludeShortAndExtendedAddress = ReadDataFlag(0x04)
)

type SecurityMode byte

const (
	// No security
	SecurityNone = SecurityMode(iota)
	// Prefonfigured Network Key
	SecurityPreconfiguredNK
	// Network Key from Trust Center
	SecurityNKTC
	// No master but Trust Center Link Key
	SecurityNoMasterTCLK
)

type PANIDMode byte

const (
	// Note predefined; The NWK PANID will be selected or obtained dynamically
	PANIDModeNotPredefined = PANIDMode(iota)
	// Predefined; The value of parameter NWK PANID (0x05) will be used to join or form a network
	PANIDModePredefined
)

type MacCapabilities uint8

const (
	// Set if this node has the capability of becoming a PAN Coordinator
	MacCapAltCoord = MacCapabilities(1 << iota)
	// Set if this node is a Full Functioning Device (FFD)
	MacCapFFD
	// Set if this node is powered by mains
	MacCapPowerSrc
	// Set if this node does not turn off the receiver when idle to conserve battery
	MacCapReceiverWhenIdle

	// Set if the device is capable of sending and receiving frames secured using the security suite specified in [B1] ???
	MacCapSecurity = MacCapabilities(1<<iota + 6)
	// The allocate address sub-field is one bit in length and shall be set to 0 or 1
	MacCapAllocAddr
)

type GPFrameType uint8

const (
	GPFrameTypeData        = GPFrameType(0b00)
	GPFrameTypeMaintenance = GPFrameType(0b01)
)

type updateNeighborAction uint8

const (
	actionRemove = updateNeighborAction(0x00)
	actionAdd    = updateNeighborAction(0x01)
)

type sendDataFlags uint8

const (
	sendDataFlagSourceRouting = sendDataFlags(0x02)
)
