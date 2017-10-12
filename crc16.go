// See http://www.ross.net/crc/download/crc_v3.txt for Parameterized CRC model,
// and http://reveng.sourceforge.net/crc-catalogue/16.htm for the CRC catalogue.
package crc16

import (
	"encoding/binary"
	"sync"
)

// A Config specifies how to calculate crc values.
type Config struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16

	// For Bytes()
	ByteOrder binary.ByteOrder
}

// A Table has its Config and room for precalculated CRC values.
type Table struct {
	Config
	tab  [256]uint16
	once sync.Once
}

// For Bytes()
var (
	BigEndian    = binary.BigEndian
	LittleEndian = binary.LittleEndian
)

// Predefined CRC16 specifications
var (
	XModem = Config{
		Poly:      0x1021,
		ByteOrder: BigEndian,
	}

	Kermit = Config{
		Poly:      0x1021,
		RefIn:     true,
		RefOut:    true,
		ByteOrder: LittleEndian,
	}

	CCITTFalse = Config{
		Poly:      0x1021,
		Init:      0xFFFF,
		ByteOrder: BigEndian,
	}

	Modbus = Config{
		Poly:      0x8005,
		Init:      0xFFFF,
		RefIn:     true,
		RefOut:    true,
		ByteOrder: LittleEndian,
	}
)

func New(c Config) *Table {
	t := &Table{Config: c}
	return t
}

func reflect8(v uint8) uint8 {
	var r uint8
	for i := uint(0); i < 8; i++ {
		if v&1 == 1 {
			r |= 1 << (7 - i)
		}
		v >>= 1
	}
	return r
}

func reflect16(v uint16) uint16 {
	var r uint16
	for i := uint(0); i < 16; i++ {
		if v&1 == 1 {
			r |= 1 << (15 - i)
		}
		v >>= 1
	}
	return r
}

func makeTable(t *Table) {
	for i := uint16(0); i < 256; i++ {
		crc := i << 8
		for j := 0; j < 8; j++ {
			if crc&(1<<15) != 0 {
				crc = (crc << 1) ^ t.Poly
			} else {
				crc <<= 1
			}
		}
		t.tab[i] = crc
	}
}

func makeReflectedTable(t *Table) {
	for i := 0; i < 256; i++ {
		crc := uint16(reflect8(uint8(i))) << 8
		for j := 0; j < 8; j++ {
			if crc&(1<<15) != 0 {
				crc = (crc << 1) ^ t.Poly
			} else {
				crc <<= 1
			}
		}
		t.tab[i] = reflect16(crc)
	}
}

// Checksum calculates CRC value for data. It uses tab's Config and calculates
// table values for the first time.
func Checksum(data []byte, tab *Table) uint16 {
	crc := tab.Init

	if tab.RefIn {
		crc = updateReflected(crc, tab, data)
	} else {
		crc = update(crc, tab, data)
	}

	if tab.RefOut {
		crc = reflect16(crc)
	}

	return crc ^ tab.XorOut
}

// Update returns the result of adding the bytes in p to the crc.
func Update(crc uint16, tab *Table, p []byte) uint16 {
	crc ^= tab.XorOut

	if tab.RefOut {
		crc = reflect16(crc)
	}

	if tab.RefIn {
		crc = updateReflected(crc, tab, p)
	} else {
		crc = update(crc, tab, p)
	}

	if tab.RefOut {
		crc = reflect16(crc)
	}

	return crc ^ tab.XorOut
}

func update(crc uint16, tab *Table, p []byte) uint16 {
	tab.once.Do(func() { makeTable(tab) })

	for _, v := range p {
		crc = tab.tab[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

func updateReflected(crc uint16, tab *Table, p []byte) uint16 {
	tab.once.Do(func() { makeReflectedTable(tab) })

	crc = reflect16(crc)
	for _, v := range p {
		crc = tab.tab[byte(crc)^v] ^ (crc >> 8)
	}
	return reflect16(crc)
}

// Bytes returns the crc value with a byte slice. It uses Table's ByteOrder. If
// ByteOrder is nil, it returns nil.
func Bytes(crc uint16, tab *Table) []byte {
	if tab.ByteOrder == nil {
		return nil
	}

	b := make([]byte, 2)
	tab.ByteOrder.PutUint16(b, crc)

	return b
}
