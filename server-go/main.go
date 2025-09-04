package main

import (
	"bufio"
	"cmp"
	"database/sql"
	"encoding/binary"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"io"
	"log"
	"net"
	"server/network"
	"server/network/c2s"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
	"server/util"
	"sync"
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
		TrackedObjs: make(map[util.Vector2I]*shared.GameObj),
		TrackedNpcs: make(map[util.Vector2I]*shared.GameNpc),
		Mu:          sync.RWMutex{},
		MaxX:        0,
		MaxY:        0,
		CurrentTick: 0,
	}
	assetsDirectory = "../../grpg-assets/"
	scriptManager   *scripts.ScriptManager
)

type ChanPacket struct {
	Bytes      []byte
	Player     *shared.Player
	PacketData c2s.PacketData
}

func main() {
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

	objs, err := LoadObjs(assetsDirectory + "objs.grpgobj")
	if err != nil {
		log.Fatal("Failed loading objs: ", err)
	}

	items, err := LoadItems(assetsDirectory + "items.grpgitem")
	if err != nil {
		log.Fatal("Failed loading items: ", err)
	}

	npcs, err := LoadNpcs(assetsDirectory + "npcs.grpgnpc")
	if err != nil {
		log.Fatal("Failed loading npcs: ", err)
	}

	LoadMaps(assetsDirectory+"maps/", g, objs)

	scriptManager = scripts.NewScriptManager()

	scriptManager.LoadObjConstants(objs)
	scriptManager.LoadItemConstants(items)
	scriptManager.LoadNpcConstants(npcs)

	err = scriptManager.LoadScripts("../game-scripts", g, npcs)
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
				packet.PacketData.Handler.Handle(buf, g, packet.Player, scriptManager)
			default:
				break processPackets
			}
		}

		timed, ok := scriptManager.TimedScripts[g.CurrentTick]
		if ok {
			for _, script := range timed {
				evaluator.Eval(script.Script, script.Env)
			}
		}

		g.CurrentTick++
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
		if err != nil {
			log.Printf("Failed to read packet opcode: %v, Conn lost.\n", err)

			game.Mu.RLock()
			player, exists := game.Connections[conn]
			game.Mu.RUnlock()

			if !exists {
				log.Printf("Couldn't find player to remove after losing connection.")
			} else {
				player.SaveToDB(game.Database)
				game.Mu.Lock()
				delete(game.Players, player)
				game.Mu.Unlock()
				network.UpdatePlayersByChunk(player.ChunkPos, game, &s2c.PlayersUpdate{ChunkPos: player.ChunkPos})
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

		game.Mu.RLock()
		player := game.Connections[conn]
		game.Mu.RUnlock()

		packets <- ChanPacket{
			Bytes:      bytes,
			Player:     player,
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
		Pos:       zeroPos,
		ChunkPos:  zeroPos,
		Facing:    shared.UP,
		Name:      string(name),
		Inventory: shared.Inventory{Items: [24]shared.InventoryItem{}},
		DialogueQueue: shared.DialogueQueue{
			Dialogues: []shared.Dialogue{},
		},
		Conn: conn,
	}

	err := player.LoadFromDB(game.Database)
	if err != nil {
		log.Printf("failed to load existing player from db, sending login rejected, err: %v\n", err)
		network.SendPacket(conn, &s2c.LoginRejected{}, game)
		return
	}

	game.Mu.Lock()
	game.Players[player] = struct{}{}
	game.Connections[conn] = player
	game.Mu.Unlock()

	network.SendPacket(conn, &s2c.LoginAccepted{}, game)
	network.UpdatePlayersByChunk(player.ChunkPos, game, &s2c.PlayersUpdate{ChunkPos: player.ChunkPos})
	network.SendPacket(player.Conn, &s2c.ObjUpdate{ChunkPos: player.ChunkPos, Rebuild: true}, game)
	network.SendPacket(player.Conn, &s2c.NpcUpdate{ChunkPos: player.ChunkPos}, game)
	network.SendPacket(player.Conn, &s2c.InventoryUpdate{Player: player}, game)
}
