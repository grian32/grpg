package shared

import (
	"database/sql"
	"net"
	"server/util"
)

type Player struct {
	Pos util.Vector2I
	// might not need these will see how design pans out
	ChunkPos  util.Vector2I
	Facing    Direction
	Inventory [24]InventoryItem
	Name      string
	Conn      net.Conn
}

func (p *Player) LoadFromDB(db *sql.DB) error {
	row := db.QueryRow("SELECT x, y FROM players WHERE name = ?", p.Name)

	var loadedX int
	var loadedY int
	err := row.Scan(&loadedX, &loadedY)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	pos := util.Vector2I{X: uint32(loadedX), Y: uint32(loadedY)}
	chunkPos := util.Vector2I{X: uint32(loadedX / 16), Y: uint32(loadedY / 16)}

	p.Pos = pos
	p.ChunkPos = chunkPos

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
		stmt, err := tx.Prepare("INSERT INTO players(player_id, name, x, y) VALUES (NULL, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.Name, p.Pos.X, p.Pos.Y)
		if err != nil {
			return err
		}
	} else {
		stmt, err := tx.Prepare("UPDATE players SET x=?, y=? WHERE player_id=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.Pos.X, p.Pos.Y, existingId)
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
