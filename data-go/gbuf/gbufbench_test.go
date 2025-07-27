package gbuf

import "testing"

func BenchmarkAll(b *testing.B) {
	for b.Loop() {
		buf := NewEmptyGBuf()

		buf.WriteByte(0x00)
		buf.WriteBytes([]byte{0x00, 0x02})
		buf.WriteInt32(-42)
		buf.WriteString("Hello, World!")
		buf.WriteUint16(42)
		buf.WriteUint32(52)

		_, _ = buf.ReadUint32()
		_, _ = buf.ReadUint16()
		_, _ = buf.ReadString()
		_, _ = buf.ReadInt32()
		_, _ = buf.ReadBytes(2)
		_, _ = buf.ReadByte()
	}
}

func BenchmarkString(b *testing.B) {
	str := "Hello, World!"
	buf := NewEmptyGBuf()

	b.ResetTimer()

	for b.Loop() {
		buf.WriteString(str)
		_, _ = buf.ReadString()
		buf.Clear()
	}
}

func BenchmarkLongString(b *testing.B) {
	str := ""
	for _ = range 128 {
		str += "A"
	}
	buf := NewEmptyGBuf()

	b.ResetTimer()

	for b.Loop() {
		buf.WriteString(str)
		_, _ = buf.ReadString()
		buf.Clear()
	}
}
