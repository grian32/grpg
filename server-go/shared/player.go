package shared

import (
	"database/sql"
	"log"
	"net"
	"server/util"
)

type Player struct {
	Pos util.Vector2I
	// might not need these will see how design pans out
	ChunkPos      util.Vector2I
	Facing        Direction
	Inventory     Inventory
	Name          string
	DialogueQueue DialogueQueue
	Conn          net.Conn
}

func (p *Player) LoadFromDB(db *sql.DB) error {
	row := db.QueryRow("SELECT x, y, inventory FROM players WHERE name = ?", p.Name)

	var loadedX int
	var loadedY int
	var invBlob []byte
	err := row.Scan(&loadedX, &loadedY, &invBlob)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	pos := util.Vector2I{X: uint32(loadedX), Y: uint32(loadedY)}
	chunkPos := util.Vector2I{X: uint32(loadedX / 16), Y: uint32(loadedY / 16)}
	inv, err := DecodeInventoryFromBlob(invBlob)
	if err != nil {
		return err
	}

	p.Pos = pos
	p.ChunkPos = chunkPos
	p.Inventory = inv

	return nil
}

func (p *Player) SaveToDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT player_id FROM players WHERE name = ?", p.Name)
	var existingId int
	err = row.Scan(&existingId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := tx.Prepare("INSERT INTO players(player_id, name, x, y, inventory) VALUES (NULL, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.Name, p.Pos.X, p.Pos.Y, p.Inventory.EncodeToBlob())
		if err != nil {
			return err
		}
	} else {
		stmt, err := tx.Prepare("UPDATE players SET x=?, y=?, inventory=? WHERE player_id=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.Pos.X, p.Pos.Y, p.Inventory.EncodeToBlob(), existingId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) GetFacingCoord() util.Vector2I {
	switch p.Facing {
	case UP:
		return util.Vector2I{X: p.Pos.X, Y: p.Pos.Y + 1}
	case RIGHT:
		return util.Vector2I{X: p.Pos.X - 1, Y: p.Pos.Y}
	case DOWN:
		return util.Vector2I{X: p.Pos.X + 1, Y: p.Pos.Y}
	case LEFT:
		return util.Vector2I{X: p.Pos.X, Y: p.Pos.Y - 1}
	default:
		log.Fatalf("unexpected Direction: %#v", p.Facing)
	}
	return util.Vector2I{}
}
