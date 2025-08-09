package grpgobj

import (
	"cmp"
	"grpg/data-go/gbuf"
	"slices"
)

type Header struct {
	Magic [8]byte
}

type ObjFlag byte

type ObjFlags byte

const (
	STATE    ObjFlag = 1 << iota // bit 0
	INTERACT                     // bit 1
)

// way this is meant to work is that the idx of Textures serves as the state number, so state 0 = x texture, state 1 = y texture, etc, each state gotta have a texture.
type Obj struct {
	Name         string
	ObjId        uint16
	Flags        ObjFlags
	Textures     []uint16 // only size 1 if non stateful
	InteractText string   // only filled in if is interact
}

// Equal only meant to be used for testing, you probably shouldn't be comparing this type
func (o *Obj) Equal(other Obj) bool {
	return o.Name == other.Name && o.ObjId == other.ObjId && o.Flags == other.Flags && slices.Equal(o.Textures, other.Textures)
}

func WriteHeader(buf *gbuf.GBuf) {
	buf.WriteBytes([]byte("GRPGOBJ\x00"))
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	bytes, err := buf.ReadBytes(8)
	if err != nil {
		return Header{}, err
	}

	return Header{
		Magic: [8]byte(bytes),
	}, nil
}

func WriteObjs(buf *gbuf.GBuf, objs []Obj) {
	buf.WriteUint16(uint16(len(objs)))

	for _, obj := range objs {
		buf.WriteString(obj.Name)
		buf.WriteUint16(obj.ObjId)
		buf.WriteByte(byte(obj.Flags))

		if !IsFlagSet(obj.Flags, STATE) {
			buf.WriteUint16(obj.Textures[0])
		} else {
			buf.WriteUint16(uint16(len(obj.Textures)))
			for _, tex := range obj.Textures {
				buf.WriteUint16(tex)
			}
		}

		if IsFlagSet(obj.Flags, INTERACT) {
			buf.WriteString(obj.InteractText)
		}
	}
}

func ReadObjs(buf *gbuf.GBuf) ([]Obj, error) {
	len, err := buf.ReadUint16()

	if err != nil {
		return nil, err
	}

	objArr := make([]Obj, len)

	for idx := range len {
		name, err1 := buf.ReadString()
		objId, err2 := buf.ReadUint16()
		flagByte, err3 := buf.ReadByte()

		if err := cmp.Or(err1, err2, err3); err != nil {
			return nil, err
		}

		// prealloc for non stateful objs, required to have atleast one tex per obj anyway
		textures := make([]uint16, 0, 1)
		flags := ObjFlags(flagByte)

		if !IsFlagSet(flags, STATE) {
			texId, err := buf.ReadUint16()

			if err != nil {
				return nil, err
			}
			textures = append(textures, texId)
		} else {
			texLen, err := buf.ReadUint16()
			if err != nil {
				return nil, err
			}

			for _ = range texLen {
				texId, err := buf.ReadUint16()
				if err != nil {
					return nil, err
				}
				textures = append(textures, texId)
			}
		}

		interactText := ""

		if IsFlagSet(flags, INTERACT) {
			interactText, err = buf.ReadString()
			if err != nil {
				return nil, err
			}
		}

		objArr[idx] = Obj{
			Name:         name,
			ObjId:        objId,
			Flags:        flags,
			Textures:     textures,
			InteractText: interactText,
		}
	}

	return objArr, nil
}

func IsFlagSet(flags ObjFlags, flag ObjFlag) bool {
	return flags&ObjFlags(flag) != 0
}
