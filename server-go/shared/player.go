package shared

import (
	"database/sql"
	"log"
	"net"
	"server/util"
)

type SkillInfo struct {
	Level uint8
	XP    uint32
}

type Player struct {
	Pos util.Vector2I
	// might not need these will see how design pans out
	ChunkPos      util.Vector2I
	Facing        Direction
	Inventory     Inventory
	Name          string
	DialogueQueue DialogueQueue
	Skills		  map[Skill]*SkillInfo
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

	// TODO
	p.Skills = make(map[Skill]*SkillInfo)
	p.Skills[FORAGING] = &SkillInfo{Level: 1, XP: 0}

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
		return util.Vector2I{X: p.Pos.X, Y: p.Pos.Y - 1}
	case RIGHT:
		return util.Vector2I{X: p.Pos.X + 1, Y: p.Pos.Y}
	case DOWN:
		return util.Vector2I{X: p.Pos.X, Y: p.Pos.Y + 1}
	case LEFT:
		return util.Vector2I{X: p.Pos.X - 1, Y: p.Pos.Y}
	default:
		log.Fatalf("unexpected Direction: %#v", p.Facing)
	}
	return util.Vector2I{}
}

func (p *Player) AddXp(skill Skill, xpAmount uint32) {
	xp := p.Skills[skill].XP
	if xp >= util.MAX_XP {
		return
	}

	if xp + xpAmount >= util.MAX_XP {
		p.Skills[skill].XP = util.MAX_XP;
		return
	}

	p.Skills[skill].XP = xp + xpAmount

	newXp := p.Skills[skill].XP

	if p.Skills[skill].Level < 75 {
		for i := p.Skills[skill].Level; i < 74; i++ {
			if newXp > util.LEVEL_XP[i] {
				p.Skills[skill].Level = uint8(i + 1);
				break;
			}
		}
	}
}
