package frame

import (
	"encoding/hex"
	slip "github.com/daedaluz/goslip"
	"testing"
)

var msg = "C01202001D00160024010400024889010401060001050010280B0A0000005EFEC0"

func TestCRC(t *testing.T) {
	slipData, _ := hex.DecodeString(msg)
	data, _ := slip.DecodeAllFromBytes(slipData)
	f := Frame(data[1])
	if !f.CheckCRC() {
		t.Fatal("CRC failed, expected:", f.getCRC(), "got", f)
	}
}
