package crc16_test

import (
	"fmt"

	"github.com/kim89098/crc16"
)

func ExampleChecksum() {
	c := crc16.New(crc16.Kermit)
	crc := crc16.Checksum([]byte("123456789"), c)
	fmt.Printf("CRC value is %#v\n", crc)
	// Output: CRC value is 0x2189
}

func ExampleUpdate() {
	c := crc16.New(crc16.Kermit)
	crc := crc16.Checksum([]byte("123456789"), c)

	if u := crc16.Update(crc, c, crc16.Bytes(crc, c)); u != 0 {
		fmt.Println("Data is corrupted")
	} else {
		fmt.Println("Data is correct")
	}
	// Output: Data is correct
}

func ExampleConfig_x25() {
	transmitter := crc16.New(crc16.Config{
		Poly:      0x1021,
		Init:      0xffff,
		RefIn:     true,
		RefOut:    true,
		XorOut:    0xffff,
		ByteOrder: crc16.LittleEndian,
	})

	b := []byte("123456789")
	crc := crc16.Checksum(b, transmitter)

	// Append CRC to data
	b = append(b, crc16.Bytes(crc, transmitter)...)

	receiver := crc16.New(crc16.Config{
		Poly:   0x1021,
		Init:   0xffff,
		RefIn:  true,
		RefOut: true,
	})

	if crc16.Checksum(b, receiver) != 0xf0b8 {
		fmt.Println("Data is corrupted")
	} else {
		fmt.Println("Data is correct")
	}
	// Output: Data is correct
}

func ExampleBytes() {
	// It uses little endian
	k := crc16.New(crc16.Kermit)
	kval := crc16.Checksum([]byte("123456789"), k)
	fmt.Printf("%#v %#v\n", kval, crc16.Bytes(kval, k))

	// It uses big endian
	x := crc16.New(crc16.XModem)
	xval := crc16.Checksum([]byte("123456789"), x)
	fmt.Printf("%#v %#v\n", xval, crc16.Bytes(xval, x))

	// Something not assigned byte order
	s := crc16.New(crc16.Config{Poly: 0x1021})
	sval := crc16.Checksum([]byte("123456789"), s)
	fmt.Printf("%#v %#v\n", sval, crc16.Bytes(sval, s))

	// Output:
	// 0x2189 []byte{0x89, 0x21}
	// 0x31c3 []byte{0x31, 0xc3}
	// 0x31c3 []byte(nil)
}
