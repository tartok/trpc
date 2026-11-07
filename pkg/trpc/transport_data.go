package trpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type TransportData struct {
	Index int16
	Size  int32
	Crc1  byte
	Data  []byte
	Crc2  byte
	r     io.Reader
	w     io.Writer
}

func NewTransportReader(r io.Reader) *TransportData {
	return &TransportData{r: r}
}

func Crc(data []byte) byte {
	var c byte
	for _, b := range data {
		c ^= b
	}
	return c
}

func ReadData(r io.Reader) ([]byte, int, error) {
	header := make([]byte, 7)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, -1, err
	}

	index := int(binary.BigEndian.Uint16(header[0:2]))
	size := int(binary.BigEndian.Uint32(header[2:6]))
	crc1 := header[6]

	if Crc(header[0:6]) != crc1 {
		return nil, -1, errors.New("crc1 error")
	}

	if size < 0 {
		return nil, index, errors.New("invalid data size")
	}

	data := make([]byte, size)
	if size > 0 {
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, index, err
		}
	}

	var crc2 [1]byte
	if _, err := io.ReadFull(r, crc2[:]); err != nil {
		return nil, index, err
	}

	if Crc(data) != crc2[0] {
		return nil, index, errors.New("crc2 error")
	}

	return data, index, nil
}

func WriteData(w io.Writer, index int, data []byte) error {
	buff := bytes.NewBuffer(make([]byte, 0, len(data)+8))

	err := binary.Write(buff, binary.BigEndian, int16(index))
	if err != nil {
		return err
	}
	err = binary.Write(buff, binary.BigEndian, int32(len(data)))
	if err != nil {
		return err
	}
	crc := Crc(buff.Bytes())
	err = binary.Write(buff, binary.BigEndian, crc)
	if err != nil {
		return err
	}
	_, err = buff.Write(data)
	if err != nil {
		return err
	}
	err = binary.Write(buff, binary.BigEndian, Crc(data))
	if err != nil {
		return err
	}
	_, err = w.Write(buff.Bytes())
	return err
}
