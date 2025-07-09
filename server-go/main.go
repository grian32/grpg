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
	"server/network"
	"server/network/c2s"
	"server/network/s2c"
	"server/shared"
	"server/util"
	"time"
)

var (
	g = &shared.Game{
		Players: []*shared.Player{},
		MaxX:    0,
		MaxY:    0,
	}
)

type ChanPacket struct {
	Bytes      []byte
	PlayerPos  int
	PacketData c2s.PacketData
}

func main() {
	LoadCollisionMaps(g)

	listener, err := net.Listen("tcp", ":4422")
	if err != nil {
		log.Fatal("Failed to start: ", err)
	}

	packets := make(chan ChanPacket, 1000)

	go cycle(packets)

	defer listener.Close()
	log.Println("Listening on 127.0.0.1:4422")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleClient(conn, g, packets)
	}
}

func cycle(packets chan ChanPacket) {
	for {
		expectedTime := time.Now().Add(60 * time.Millisecond)

		// dodgy label hack to break out properly but wcyd
	processPackets:
		for {
			select {
			case packet := <-packets:
				buf := gbuf.NewGBuf(packet.Bytes)
				packet.PacketData.Handler.Handle(buf, g, packet.PlayerPos)
			default:
				break processPackets
			}
		}

		diff := expectedTime.Sub(time.Now())
		if diff > 0 {
			time.Sleep(diff)
		}
	}
}

func handleClient(conn net.Conn, game *shared.Game, packets chan ChanPacket) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	log.Printf("Client connected with ip %s\n", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		opcode, err := reader.ReadByte()
		fmt.Println(opcode)
		if err != nil {
			log.Printf("Failed to read packet opcode: %v, Conn lost.\n", err)
			playerPos := -1
			var playerChunk util.Vector2I
			for idx, p := range game.Players {
				if p.Conn == conn {
					playerPos = idx
					playerChunk = p.ChunkPos
					break
				}
			}

			if playerPos == -1 {
				log.Printf("Couldn't find player to remove after losing connection.")
			} else {
				game.Players[playerPos] = game.Players[len(game.Players)-1]
				game.Players = game.Players[:len(game.Players)-1]
				network.UpdatePlayersByChunk(playerChunk, game)
			}

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

		var playerPos = -1

		for idx, p := range game.Players {
			if p.Conn == conn {
				playerPos = idx
			}
		}

		if playerPos == -1 {
			fmt.Printf("Couldn't find player with conn %v", conn)
			return
		}

		//packetData.Handler.Handle(buf, game, playerPos)
		packets <- ChanPacket{
			Bytes:      bytes,
			PlayerPos:  playerPos,
			PacketData: packetData,
		}
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
			network.SendPacket(conn, &s2c.LoginRejected{}, game)
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

	network.SendPacket(conn, &s2c.LoginAccepted{}, game)
	// this will be changed to the chunkpos where u login when i have player saves
	network.UpdatePlayersByChunk(util.Vector2I{X: 0, Y: 0}, game)
}
