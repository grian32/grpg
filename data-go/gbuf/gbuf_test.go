package gbuf

import (
	"bytes"
	"cmp"
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

func TestGBuf_WriteUint16(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteUint16(0x1234)

	expected := []byte{0x12, 0x34}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteUint16() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteUint16Multiple(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteUint16(0x1234)
	buf.WriteUint16(0x5678)

	expected := []byte{0x12, 0x34, 0x56, 0x78}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteUint16() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteUint32(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteUint32(0x12345678)

	expected := []byte{0x12, 0x34, 0x56, 0x78}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteUint32() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteUint32Multiple(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteUint32(0x12345678)
	buf.WriteUint32(0x9ABCDEF0)

	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteUint32() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteBytes(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteBytes([]byte{0x01, 0x02, 0x03})

	expected := []byte{0x01, 0x02, 0x03}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteBytes() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteBytesMultiple(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteBytes([]byte{0x01, 0x02})
	buf.WriteBytes([]byte{0x03, 0x04, 0x05})

	expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteBytes() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteBytesEmpty(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteBytes([]byte{})

	expected := []byte{}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteBytes() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteString(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteString("hello")

	// Length 5 (0x00000005) followed by "hello"
	expected := []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteString() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteStringEmpty(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteString("")

	// Length 0 (0x00000000) followed by no bytes
	expected := []byte{0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteString() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_WriteStringMultiple(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteString("hi")
	buf.WriteString("bye")

	// Length 2 (0x00000002) + "hi" + Length 3 (0x00000003) + "bye"
	expected := []byte{0x00, 0x00, 0x00, 0x02, 'h', 'i', 0x00, 0x00, 0x00, 0x03, 'b', 'y', 'e'}
	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_WriteString() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_Clear(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewGBuf(data)
	buf.pos = 2

	buf.Clear()

	expectedSlice := []byte{}
	expectedPos := 0

	if !bytes.Equal(expectedSlice, buf.slice) || buf.pos != expectedPos {
		t.Errorf("GBuf_Clear() slice=%v, pos=%d want slice=%v, pos=%d", buf.slice, buf.pos, expectedSlice, expectedPos)
	}
}

func TestGBuf_ClearEmptyBuffer(t *testing.T) {
	buf := NewGBuf([]byte{})
	buf.pos = 0

	buf.Clear()

	expectedSlice := []byte{}
	expectedPos := 0

	if !bytes.Equal(expectedSlice, buf.slice) || buf.pos != expectedPos {
		t.Errorf("GBuf_Clear() slice=%v, pos=%d want slice=%v, pos=%d", buf.slice, buf.pos, expectedSlice, expectedPos)
	}
}

func TestGBuf_ClearThenWrite(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewGBuf(data)

	buf.Clear()
	buf.WriteUint16(0x1234)

	expected := []byte{0x12, 0x34}

	if !bytes.Equal(expected, buf.slice) {
		t.Errorf("GBuf_Clear() then WriteUint16() got %v, want %v", buf.slice, expected)
	}
}

func TestGBuf_Bytes(t *testing.T) {
	buf := NewEmptyGBuf()
	buf.WriteUint32(2)

	expected := []byte{0x00, 0x00, 0x00, 0x02}

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("GBuf_Bytes()=%v want match for %v", buf.Bytes(), expected)
	}
}

func TestGBuf_ReadByte(t *testing.T) {
	buf := NewGBuf([]byte{0x00, 0x01})

	firstByte, err1 := buf.ReadByte()
	secondByte, err2 := buf.ReadByte()

	err := cmp.Or(err1, err2)

	if firstByte != 0x00 || secondByte != 0x01 || err != nil {
		t.Errorf("GBuf_ReadByte()=%b, %b & %v want match for %b, %b", firstByte, secondByte, err, 0x00, 0x01)
	}
}

func TestGBuf_WriteByte(t *testing.T) {
	buf := NewEmptyGBuf()

	buf.WriteByte(0xFF)
	buf.WriteByte(0x01)
	buf.WriteByte(0x02)

	expectedBytes := []byte{0xFF, 0x01, 0x02}

	if !bytes.Equal(buf.Bytes(), expectedBytes) {
		t.Errorf("GBuf_ReadByte()=%v want match for %v", buf.Bytes(), expectedBytes)
	}
}
