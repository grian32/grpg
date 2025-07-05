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

func NewEmptyGBuf() *GBuf {
	return &GBuf{
		slice: []byte{},
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

func (buf *GBuf) ReadByte() (byte, error) {
	if buf.pos+1 > len(buf.slice) {
		return 0x00, errors.New("not enough byte to read bytes")
	}
	val := buf.slice[buf.pos]
	buf.pos += 1
	return val, nil
}

// TODO: tests
func (buf *GBuf) ReadInt32() (int32, error) {
	if buf.pos+4 > len(buf.slice) {
		return 0, errors.New("not enough bytes to read int32")
	}
	val := int32(binary.BigEndian.Uint32(buf.slice[buf.pos : buf.pos+4]))
	buf.pos += 4
	return val, nil
}

func (buf *GBuf) WriteUint16(val uint16) {
	temp := make([]byte, 2)
	binary.BigEndian.PutUint16(temp, val)
	buf.slice = append(buf.slice, temp...)
}

func (buf *GBuf) WriteUint32(val uint32) {
	temp := make([]byte, 4)
	binary.BigEndian.PutUint32(temp, val)
	buf.slice = append(buf.slice, temp...)
}

func (buf *GBuf) WriteBytes(bytes []byte) {
	buf.slice = append(buf.slice, bytes...)
}

// WriteString
// Writes a uint32 length encoded string to the buffer
func (buf *GBuf) WriteString(val string) {
	buf.WriteUint32(uint32(len(val)))
	buf.WriteBytes([]byte(val))
}

func (buf *GBuf) WriteByte(val byte) {
	buf.slice = append(buf.slice, val)
}

// Clear
// Replaces the backing slice with an empty one
func (buf *GBuf) Clear() {
	buf.slice = make([]byte, 0)
	buf.pos = 0
}

func (buf *GBuf) Bytes() []byte {
	return buf.slice
}
