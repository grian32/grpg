package grpgobj

import (
	"bytes"
	"grpg/data-go/gbuf"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedHeader := Header{
		Magic: [8]byte{'G', 'R', 'P', 'G', 'O', 'B', 'J', 0x00},
	}

	expectedBytes := [8]byte{'G', 'R', 'P', 'G', 'O', 'B', 'J', 0x00}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteHeader", func(t *testing.T) {
		WriteHeader(buf)

		if !bytes.Equal(buf.Bytes(), expectedBytes[:]) {
			t.Errorf("WriteHeader=%v, want=%v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadHeader", func(t *testing.T) {
		header, err := ReadHeader(buf)

		if header != expectedHeader || err != nil {
			t.Errorf("ReadHeader=%v,%v, want=%v", header, err, expectedHeader)
		}
	})
}

func TestReadWriteObjs(t *testing.T) {
	expectedObjs := []Obj{
		{
			Name:     "stone",
			ObjId:    2,
			Flags:    0,
			Textures: []uint16{1},
		},
		{
			Name:     "berry_bush",
			ObjId:    4,
			Flags:    ObjFlags(STATE | INTERACT),
			Textures: []uint16{2, 3},
		},
	}

	expectedBytes := []byte{
		0x00, 0x02, // count

		0x00, 0x00, 0x00, 0x05, 's', 't', 'o', 'n', 'e', // name
		0x00, 0x02, // obj id
		0x00,       // flags
		0x00, 0x01, // tex id since STATE is not set

		0x00, 0x00, 0x00, 0x0A, 'b', 'e', 'r', 'r', 'y', '_', 'b', 'u', 's', 'h', // name
		0x00, 0x04, // obj id
		0x03,       // flag, 00 00 00 11
		0x00, 0x02, // tex len
		0x00, 0x02, 0x00, 0x03, // tex arr
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteObjs", func(t *testing.T) {
		WriteObjs(buf, expectedObjs)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteObjs=%v, want=%v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadObjs", func(t *testing.T) {
		objs, err := ReadObjs(buf)

		if len(objs) != 2 || !objs[0].Equal(expectedObjs[0]) || !objs[1].Equal(expectedObjs[1]) || err != nil {
			t.Errorf("ReadObjs=%v,%v, want=%v", objs, err, expectedObjs)
		}
	})
}

func TestIsFlagSet(t *testing.T) {
	if !IsFlagSet(ObjFlags(STATE|INTERACT), STATE) {
		t.Errorf("expected STATE flag to be set, but it was not")
	}

	if IsFlagSet(ObjFlags(INTERACT), STATE) {
		t.Errorf("expected STATE flag to not be set, but it was")
	}

	if IsFlagSet(ObjFlags(0), STATE) {
		t.Error("STATE flag is set in empty ObjFlags")
	}
}
