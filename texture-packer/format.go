package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"log"
	"os"
)

type GRPGTexHeader struct {
	Magic   [8]byte
	Version uint16
}

type GRPGTexTexture struct {
	InternalIdLength uint16
	InternalIdData   []byte
	RGBPixels        [12288]byte // 64 px * 64 px * 3 bytes(r, g, b)
}

func WriteGRPGTexHeader(buf *bytes.Buffer, version uint16) {
	header := GRPGTexHeader{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: version,
	}
	err := binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		log.Fatal(err)
	}
}

func BuildGRPGTexFromManifest(files []GRPGTexManifestEntry) []GRPGTexTexture {
	tex := make([]GRPGTexTexture, len(files))

	for idx, file := range files {
		f, err := os.Open(file.FilePath)
		if err != nil {
			log.Fatal(err)
		}

		data, err := png.Decode(f)

		if err != nil {
			log.Fatal(err)
		}

		maxBounds := data.Bounds().Max

		if maxBounds.X != 64 || maxBounds.Y != 64 {
			log.Fatal("textures that are not exactly 64x64 are disallowed")
		}

		rgbArray := [12288]byte{}

		rgbIdx := 0

		for y := range 64 {
			for x := range 64 {
				r, g, b, _ := data.At(x, y).RGBA()
				rgbArray[rgbIdx] = byte(r >> 8)
				rgbArray[rgbIdx+1] = byte(g >> 8)
				rgbArray[rgbIdx+2] = byte(b >> 8)
				rgbIdx += 3
			}
		}

		tex[idx] = GRPGTexTexture{
			InternalIdLength: uint16(len(file.InternalName)),
			InternalIdData:   []byte(file.InternalName),
			RGBPixels:        rgbArray,
		}

		f.Close()
	}

	return tex
}

func WriteGRPGTex(buf *bytes.Buffer, textures []GRPGTexTexture) {
	err := binary.Write(buf, binary.BigEndian, uint32(len(textures)))
	if err != nil {
		log.Fatal(err)
	}

	for _, tex := range textures {
		err = binary.Write(buf, binary.BigEndian, tex.InternalIdLength)

		if err != nil {
			log.Fatal(err)
		}
		// using buf.write cuz not fixed length lol
		buf.Write(tex.InternalIdData)

		err = binary.Write(buf, binary.BigEndian, tex.RGBPixels)

		if err != nil {
			log.Fatal(err)
		}
	}
}
