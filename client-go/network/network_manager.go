package network

import (
	"bufio"
	"client/network/s2c"
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

		fmt.Println(opcode)

		packetData := s2c.Packets[opcode]

		bytes := make([]byte, packetData.Length)

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
