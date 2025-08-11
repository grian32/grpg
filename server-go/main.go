package main

import (
	"bufio"
	"cmp"
	"database/sql"
	"encoding/binary"
	"fmt"
	"grpg/data-go/gbuf"
	"io"
	"log"
	"net"
	"server/network"
	"server/network/c2s"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
	"server/util"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	g = &shared.Game{
		Players:     map[*shared.Player]struct{}{},
		Connections: make(map[net.Conn]*shared.Player),
		MaxX:        0,
		MaxY:        0,
	}
	assetsDirectory = "../../grpg-assets/"
)

type ChanPacket struct {
	Bytes      []byte
	Player     *shared.Player
	PacketData c2s.PacketData
}

func main() {
	LoadCollisionMaps(assetsDirectory+"maps/", g)

	db, err := sql.Open("sqlite3", "./players.db")
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer db.Close()

	g.Database = db

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("Failed to create sqlite3 driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
	if err != nil {
		log.Fatal("Failed to create new migrate: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to migrate: ", err)
	}

	listener, err := net.Listen("tcp", ":4422")
	if err != nil {
		log.Fatal("Failed to start: ", err)
	}

	scriptManager := scripts.ScriptManager{}
	err = scriptManager.LoadScripts("../game-scripts")
	if err != nil {
		log.Fatal("Failed loading scripts: ", err)
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
				packet.PacketData.Handler.Handle(buf, g, packet.Player)
			default:
				break processPackets
			}
		}

		diff := time.Until(expectedTime)
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

			player, exists := game.Connections[conn]

			if !exists {
				log.Printf("Couldn't find player to remove after losing connection.")
			} else {
				player.SaveToDB(game.Database)
				delete(game.Players, player)
				network.UpdatePlayersByChunk(player.ChunkPos, game)
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

		packets <- ChanPacket{
			Bytes:      bytes,
			Player:     game.Connections[conn],
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

	// TODO: not sure if i can optimize this unless i keep a map of pre existing names or something but feels overkill alrdy
	for player, _ := range game.Players {
		if player.Name == string(name) {
			network.SendPacket(conn, &s2c.LoginRejected{}, game)
			return
		}
	}

	zeroPos := util.Vector2I{X: 0, Y: 0}

	player := &shared.Player{
		Pos:      zeroPos,
		ChunkPos: zeroPos,
		Facing:   shared.UP,
		Name:     string(name),
		Conn:     conn,
	}

	game.Players[player] = struct{}{}
	game.Connections[conn] = player
	player.LoadFromDB(game.Database)

	network.SendPacket(conn, &s2c.LoginAccepted{}, game)
	// this will be changed to the chunkpos where u login when i have player saves
	network.UpdatePlayersByChunk(util.Vector2I{X: 0, Y: 0}, game)
}
