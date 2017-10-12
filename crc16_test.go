package crc16

import (
	"testing"
)

type ChecksumSet struct {
	Name   string
	Config Config
	Check  uint16
}

func TestChecksum(t *testing.T) {
	s := []ChecksumSet{
		ChecksumSet{"XModem", XModem, 0x31c3},
		ChecksumSet{"Kermit", Kermit, 0x2189},
		ChecksumSet{"CCITT-FALSE", CCITTFalse, 0x29b1},
		ChecksumSet{"Modbus", Modbus, 0x4b37},
	}

	for _, v := range s {
		c := New(v.Config)
		crc := Checksum([]byte("123456789"), c)
		if crc != v.Check {
			t.Errorf("Unexpected %v CRC value: %#v, expecting %#v", v.Name, crc, v.Check)
		}
	}
}

type UpdateSet struct {
	Name   string
	Config Config
}

func TestUpdate(t *testing.T) {
	s := []UpdateSet{
		UpdateSet{"XModem", XModem},
		UpdateSet{"Kermit", Kermit},
		UpdateSet{"CCITT-FALSE", CCITTFalse},
		UpdateSet{"Modbus", Modbus},
	}

	for _, v := range s {
		c := New(v.Config)
		crc := Checksum([]byte("123456789"), c)

		if crc := Update(crc, c, Bytes(crc, c)); crc != 0 {
			t.Errorf("Unexpected %v CRC value: %#v, expecting 0", v.Name, crc)
		}
	}
}
