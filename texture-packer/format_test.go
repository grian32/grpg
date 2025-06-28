package main

import (
	"bytes"
	"testing"
)

var buf = bytes.Buffer{}

func TestWriteGRPGTexHeaderVer1(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0x00, 0x01, // ver1
	}

	err := WriteGRPGTexHeader(&buf, 1)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteGRPGTexHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}

func TestWriteGRPGTexHeaderVerMax(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0xFF, 0xFF, // ver1
	}

	err := WriteGRPGTexHeader(&buf, 65535)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteGRPGTexHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}
