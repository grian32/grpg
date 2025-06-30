package gbuf

import (
	"encoding/binary"
	"errors"
)

type GBuf struct {
	slice []byte
	pos   int
}

func NewGBuf(data []byte) *GBuf {
	return &GBuf{
		slice: data,
		pos:   0,
	}
}

func (buf *GBuf) ReadUint16() (uint16, error) {
	if buf.pos+2 > len(buf.slice) {
		return 0, errors.New("not enough bytes to read uint16")
	}
	val := binary.BigEndian.Uint16(buf.slice[buf.pos : buf.pos+2])
	buf.pos += 2
	return val, nil
}

func (buf *GBuf) ReadUint32() (uint32, error) {
	if buf.pos+4 > len(buf.slice) {
		return 0, errors.New("not enough bytes to read uint16")
	}
	val := binary.BigEndian.Uint32(buf.slice[buf.pos : buf.pos+4])
	buf.pos += 4
	return val, nil
}

func (buf *GBuf) ReadBytes(length int) ([]byte, error) {
	if buf.pos+length > len(buf.slice) {
		return nil, errors.New("not enough bytes to read specified length of bytes")
	}
	val := buf.slice[buf.pos : buf.pos+length]
	buf.pos += length
	return val, nil
}

// ReadString
// Reads a uint32 length encoded string.
func (buf *GBuf) ReadString() (string, error) {
	length, err := buf.ReadUint32()
	if err != nil {
		return "", errors.New("failed to read uint32 length of string")
	}
	bytes, err := buf.ReadBytes(int(length))
	if err != nil {
		return "", errors.New("failed to read bytes of a string")
	}

	return string(bytes[:]), nil
}
