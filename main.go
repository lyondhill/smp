package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/sigurn/crc16"
)

type Message struct {
	ID      uint16
	Type    uint8
	Length  uint16
	Payload []byte
	CRC     uint16
}

func main() {
	proto := []byte{0x12, 0x34, 0x01, 0x00, 0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x1a, 0x2b}
	fmt.Printf("%q\n", proto)
	input, err := encode("hello")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", input)

	buf := bytes.NewBuffer(input)

	message, err := decode(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", message)

}

func encode(payload string) ([]byte, error) {
	table := crc16.MakeTable(crc16.CRC16_MAXIM)
	m := Message{
		ID:      1,
		Type:    0,
		Length:  uint16(len(payload)),
		Payload: []byte(payload),
		CRC:     crc16.Checksum([]byte(payload), table),
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.LittleEndian, m.ID); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, m.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, m.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, m.Payload); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, m.CRC); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}

func decode(reader io.Reader) (*Message, error) {
	message := &Message{}

	if err := binary.Read(reader, binary.LittleEndian, &message.ID); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &message.Type); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &message.Length); err != nil {
		return nil, err
	}

	message.Payload = make([]byte, message.Length)
	if err := binary.Read(reader, binary.LittleEndian, message.Payload); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &message.CRC); err != nil {
		return nil, err
	}

	// verify CRC
	table := crc16.MakeTable(crc16.CRC16_MAXIM)
	crc := crc16.Checksum(message.Payload, table)
	if crc != message.CRC {
		return nil, fmt.Errorf("payload CRC doesn't match checked body: %d %d", message.CRC, crc)
	}
	return message, nil
}
