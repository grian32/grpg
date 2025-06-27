package main

import (
	"bytes"
	"encoding/binary"
)

type GRPGTexHeader struct {
	Magic   [8]byte
	Version uint16
}

func WriteGRPGTexHeader(buf *bytes.Buffer, version uint16) {
	header := GRPGTexHeader{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: version,
	}
	_ = binary.Write(buf, binary.BigEndian, header)
}
