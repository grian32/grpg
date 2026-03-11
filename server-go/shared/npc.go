package shared

import (
	"grpg/data-go/grpgnpc"
	"math/rand/v2"
	"server/util"
)

type GameNpc struct {
	Pos         util.Vector2I
	Uid         uint32
	NpcData     *grpgnpc.Npc
	ChunkPos    util.Vector2I
	ValidWander map[util.Vector2I]struct{}
	WanderMin   util.Vector2I
	WanderMax   util.Vector2I
	WanderRange uint8
}

type NpcMove struct {
	From util.Vector2I
	To   util.Vector2I
}

type NpcPath struct {
	NpcId uint16
	Moves []NpcMove
}

func (g *GameNpc) Wander(game *Game) {
	if g.ValidWander == nil {
		// i think safe assumption if they dont have a valid wander is that they have never wandered and as such the curr pos is still the initial position
		startPosX := util.ClampMin(g.Pos.X-uint32(g.WanderRange), 0)
		startPosY := util.ClampMin(g.Pos.Y-uint32(g.WanderRange), 0)
		endPosX := util.ClampMax(g.Pos.X+uint32(g.WanderRange), (g.ChunkPos.X+1)*16)
		endPosY := util.ClampMax(g.Pos.Y+uint32(g.WanderRange), (g.ChunkPos.Y+1)*16)
		validPos := make(map[util.Vector2I]struct{}, 0)

		for x := startPosX; x <= endPosX; x++ {
			for y := startPosY; y <= endPosY; y++ {
				pos := util.Vector2I{X: x, Y: y}
				if _, exists := game.CollisionMap[pos]; !exists {
					validPos[pos] = struct{}{}
				}
			}
		}

		g.ValidWander = validPos
		g.WanderMin = util.Vector2I{X: startPosX, Y: startPosY}
		g.WanderMax = util.Vector2I{X: endPosX, Y: endPosY}
	}

	pos := rand.IntN(len(g.ValidWander))
	var key util.Vector2I
	i := 0
	for k := range g.ValidWander {
		if i == pos {
			key = k
			break
		}
		i++
	}
	path := BFS(g.Pos, key, g.WanderMax, g.WanderMin, g.ValidWander)
	// mby dubious extra access here? not toooo sure...
	_, exists := game.NpcMoves[g.ChunkPos]
	if !exists {
		game.NpcMoves[g.ChunkPos] = []NpcPath{{NpcId: g.NpcData.NpcId, Moves: path}}
	} else {
		game.NpcMoves[g.ChunkPos] = append(game.NpcMoves[g.ChunkPos],
			NpcPath{NpcId: g.NpcData.NpcId, Moves: path},
		)
	}

}

func BFS(start, goal, max, min util.Vector2I, valid map[util.Vector2I]struct{}) []NpcMove {
	prev := map[util.Vector2I]util.Vector2I{
		start: start,
	}
	queue := []util.Vector2I{start}
	// kind of a hack but my util.vector2i server side is uint
	dirsX := []int32{1, -1, 0, 0}
	dirsY := []int32{0, 0, 1, -1}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == goal {
			var path []NpcMove
			for p := goal; p != start; p = prev[p] {
				path = append(path, NpcMove{
					From: prev[p],
					To:   p,
				})
			}
			return path
		}

		for idx, dx := range dirsX {
			dy := dirsY[idx]
			nX := int32(curr.X) + dx
			nY := int32(curr.Y) + dy

			if nX >= int32(min.X) && nX <= int32(max.X) && nY >= int32(min.Y) && nY <= int32(max.Y) {
				n := util.Vector2I{X: uint32(nX), Y: uint32(nY)}
				_, seen := prev[n]
				_, ok := valid[n]
				if ok && !seen {
					prev[n] = curr
					queue = append(queue, n)
				}
			}
		}
	}

	return nil
}
