package main

import (
	"bufio"
	"cmp"
	"encoding/binary"
	"fmt"
	"grpg/data-go/gbuf"
	"io"
	"log"
	"net"
	"server/network/c2s"
	"server/network/s2c"
	"server/shared"
	"server/util"
)

var (
	g = &shared.Game{
		Players:        []*shared.Player{},
		PlayersByChunk: map[util.Vector2I][]*shared.Player{},
		MaxX:           15,
		MaxY:           31,
	}
)

func main() {
	listener, err := net.Listen("tcp", ":4422")
	if err != nil {
		log.Fatal("Failed to start: ", err)
	}

	defer listener.Close()
	log.Println("Listening on 127.0.0.1:4422")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleClient(conn, g)
	}
}

func handleClient(conn net.Conn, game *shared.Game) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	log.Printf("Client connected with ip %s\n", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		for _, p := range game.Players {
			fmt.Printf("%s @ %d, %d\n", p.Name, p.Pos.X, p.Pos.Y)
		}

		opcode, err := reader.ReadByte()
		if err != nil {
			log.Printf("Failed to read packet opcode: %v, Conn lost.\n", err)
			// TODO: remove player
			return
		}

		if opcode == 0x01 {
			handleLogin(reader, conn, game)
			continue
		}

		packetData := c2s.Packets[opcode]

		bytes := make([]byte, packetData.Length)

		_, err = io.ReadFull(reader, bytes)
		if err != nil {
			log.Printf("Failed to read bytes from opcode %b, %v\n", opcode, err)
			return
		}

		buf := gbuf.NewGBuf(bytes)

		var playerPos = -1

		for idx, p := range game.Players {
			if p.Conn == conn {
				playerPos = idx
			} else {
				log.Printf("Couldn't find player in position, %v\n", p)
				return
			}
		}

		packetData.Handler.Handle(buf, game, playerPos)
	}
}

func handleLogin(reader *bufio.Reader, conn net.Conn, game *shared.Game) {
	nameLenBytes := make([]byte, 4)
	_, err1 := io.ReadFull(reader, nameLenBytes)

	nameLen := binary.BigEndian.Uint32(nameLenBytes)
	name := make([]byte, nameLen)
	_, err2 := io.ReadFull(reader, name)

	if err := cmp.Or(err1, err2); err != nil {
		log.Printf("Error reading login packet, %v\n", err)
	}

	for _, player := range game.Players {
		if player.Name == string(name) {
			shared.SendPacket(conn, &s2c.LoginRejected{})
			return
		}
	}

	zeroPos := util.Vector2I{X: 0, Y: 0}

	player := &shared.Player{
		Pos:      zeroPos,
		ChunkPos: zeroPos,
		Name:     string(name),
		Conn:     conn,
	}

	game.Players = append(game.Players, player)
	game.PlayersByChunk[zeroPos] = append(game.PlayersByChunk[zeroPos], player)

	shared.SendPacket(conn, &s2c.LoginAccepted{})
}
