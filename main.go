package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

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

	if len(os.Args) < 2 {
		fmt.Println("hay bud, I need to know what im doing:")
		fmt.Println("  create = create a file every thing you type will be read, new lines denote new messages ctrl+d to close and write the file")
		fmt.Println("  read = read the file you wrote 'file.encoded'")
	}

	switch os.Args[1] {
	case "create":
		err := create()
		if err != nil {
			panic(err)
		}
		fmt.Println("stuff wrote")
	case "read":
		err := read()
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("hay dude.. what are you thinking? create or read figure it out!")
		os.Exit(1)
	}
}

func create() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("reading stuff ctrl+d to write and close")

	content := []byte{}
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Printf("%q", text)
		encoded, err := encode(text[:len(text)-1])
		if err != nil {
			return err
		}
		fmt.Printf("encoded: %x\n", encoded)
		content = append(content, encoded...)
	}

	fmt.Printf("writing stuff: %x", content)
	return os.WriteFile("file.encoded", content, 0666)
}

func read() error {
	data, err := os.ReadFile("file.encoded")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(data)

	for {
		message, err := decode(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		fmt.Printf("payload: %s \tmessage: %+v, \n", message.Payload, message)
	}
	fmt.Println("all done here!!")
	return nil
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
	if err := binary.Write(buffer, binary.BigEndian, m.ID); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, m.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.Payload); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.CRC); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}

func decode(reader io.Reader) (*Message, error) {
	message := &Message{}

	if err := binary.Read(reader, binary.BigEndian, &message.ID); err != nil {
		return nil, errors.Join(errors.New("failed to binary.Read(message.ID)"), err)
	}

	if err := binary.Read(reader, binary.BigEndian, &message.Type); err != nil {
		return nil, errors.Join(errors.New("failed to binary.Read(message.Type)"), err)
	}

	if err := binary.Read(reader, binary.BigEndian, &message.Length); err != nil {
		return nil, errors.Join(errors.New("failed to binary.Read(message.Length)"), err)
	}
	message.Payload = make([]byte, message.Length)
	if err := binary.Read(reader, binary.BigEndian, message.Payload); err != nil {
		return nil, errors.Join(errors.New("failed to binary.Read(message.Payload)"), err)
	}

	if err := binary.Read(reader, binary.BigEndian, &message.CRC); err != nil {
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
