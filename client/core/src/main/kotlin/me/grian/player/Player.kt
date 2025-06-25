package me.grian.player

import com.badlogic.gdx.Gdx
import me.grian.Main
import me.grian.network.NetworkManager
import me.grian.network.packets.c2s.C2SMovePacket

class Player(
    val pos: Point,
    var chunkPos: Point,
    var realX: Float,
    var realY: Float,
    var name: String
) {
    fun move(x: Int, y: Int) {
        // TODO: bounds check on map once i have that sorted out
        if (x !in 0..15 || y !in 0..31) return
        if (Main.players.any { x == it.pos.x && y == it.pos.y }) return

        // maybe reconsider sending packet on function call, not sure though, have to see which one
        // ends up being more used, setting x/y directly or normnal move
        NetworkManager.sendPacket(
            C2SMovePacket(x, y)
        )
    }

    constructor(x: Int, y: Int, chunkX: Int, chunkY: Int,  name: String):
        this(Point(x, y), Point(chunkX, chunkY), x * Main.tileSize, y * Main.tileSize, name)
}
