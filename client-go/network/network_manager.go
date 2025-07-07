package network

import (
	"bufio"
	"client/network/s2c"
	"encoding/binary"
	"fmt"
	"grpg/data-go/gbuf"
	"io"
	"log"
	"net"
)

type ChanPacket struct {
	Buf        *gbuf.GBuf
	PacketData s2c.PacketData
}

func StartConn() net.Conn {
	netConn, err := net.Dial("tcp", "localhost:4422")
	if err != nil {
		log.Fatal(err)
	}

	return netConn
}

func ReadServerPackets(conn net.Conn, packetChan chan<- ChanPacket) {
	reader := bufio.NewReader(conn)

	for {
		opcode, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading packet opcode, conn lost.")
			return
		}

		packetData := s2c.Packets[opcode]

		var bytes []byte

		// the kotlin one anymore, basically reads the amount of players and then reads bytes assuming 8 char player
		// names
		if packetData.Length == -1 {
			lenBytes := make([]byte, 2)
			_, err = io.ReadFull(reader, lenBytes)

			if err != nil {
				log.Printf("Failed to read packet length for variable length packet with opcode %b, %v\n", opcode, err)
				continue
			}

			packetLen := binary.BigEndian.Uint16(lenBytes)
			bytes = make([]byte, packetLen)
		} else {
			bytes = make([]byte, packetData.Length)
		}
		_, err = io.ReadFull(reader, bytes)
		if err != nil {
			return
		}

		buf := gbuf.NewGBuf(bytes)

		select {
		case packetChan <- ChanPacket{
			Buf:        buf,
			PacketData: packetData,
		}:
		default:
			return
		}
	}
}
