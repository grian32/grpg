package gbuf

import (
	"bytes"
	"testing"
)

func TestGBuf_ReadUint16(t *testing.T) {
	data := []byte{0x00, 0x02}
	buf := NewGBuf(data)

	output, err := buf.ReadUint16()
	expected := uint16(2)

	if expected != output || err != nil {
		t.Errorf("GBuf_ReadUint16()=%d, %v want match for %d", output, err, expected)
	}
}

func TestGBuf_ReadUint16NotEnoughData(t *testing.T) {
	data := []byte{0x00}
	buf := NewGBuf(data)

	output, err := buf.ReadUint16()

	if err == nil {
		t.Errorf("GBuf_ReadUint16()=%d, %v did not error", output, err)
	}
}

func TestGBuf_ReadUint32(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x04}
	buf := NewGBuf(data)

	output, err := buf.ReadUint32()
	expected := uint32(4)

	if expected != output || err != nil {
		t.Errorf("GBuf_ReadUint32()=%d, %v want match for %d", output, err, expected)
	}
}

func TestGBuf_ReadUint32NotEnoughData(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02}
	buf := NewGBuf(data)

	output, err := buf.ReadUint32()

	if err == nil {
		t.Errorf("GBuf_ReadUint32()=%d, %v did not error", output, err)
	}
}

func TestGBuf_ReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	buf := NewGBuf(data)

	output, err := buf.ReadBytes(3)
	expected := []byte{0x01, 0x02, 0x03}

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf("GBuf_ReadBytes()=%v, %v want match for %v", output, err, expected)
	}
}

func TestGBuf_ReadBytesNotEnoughData(t *testing.T) {
	data := []byte{0x01, 0x02}
	buf := NewGBuf(data)

	output, err := buf.ReadBytes(5)

	if err == nil {
		t.Errorf("GBuf_ReadBytes()=%v, %v did not error", output, err)
	}
}

func TestGBuf_ReadBytesZeroLength(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	buf := NewGBuf(data)

	output, err := buf.ReadBytes(0)
	expected := []byte{}

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf("GBuf_ReadBytes()=%v, %v want match for %v", output, err, expected)
	}
}

func TestGBuf_ReadString(t *testing.T) {
	// Length 5 (0x00000005) followed by "hello"
	data := []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}
	buf := NewGBuf(data)

	output, err := buf.ReadString()
	expected := "hello"

	if expected != output || err != nil {
		t.Errorf("GBuf_ReadString()=%s, %v want match for %s", output, err, expected)
	}
}

func TestGBuf_ReadStringNotEnoughDataForLength(t *testing.T) {
	data := []byte{0x00, 0x00, 0x01} // Incomplete uint32 length
	buf := NewGBuf(data)

	output, err := buf.ReadString()

	if err == nil {
		t.Errorf("GBuf_ReadString()=%s, %v did not error", output, err)
	}
}

func TestGBuf_ReadStringNotEnoughDataForBytes(t *testing.T) {
	// Length 10 but only 3 bytes available
	data := []byte{0x00, 0x00, 0x00, 0x0A, 'a', 'b', 'c'}
	buf := NewGBuf(data)

	output, err := buf.ReadString()

	if err == nil {
		t.Errorf("GBuf_ReadString()=%s, %v did not error", output, err)
	}
}

func TestGBuf_ReadStringEmpty(t *testing.T) {
	// Length 0 followed by no bytes
	data := []byte{0x00, 0x00, 0x00, 0x00}
	buf := NewGBuf(data)

	output, err := buf.ReadString()
	expected := ""

	if expected != output || err != nil {
		t.Errorf("GBuf_ReadString()=%s, %v want match for %s", output, err, expected)
	}
}
