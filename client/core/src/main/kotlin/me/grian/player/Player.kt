package me.grian.player

import com.badlogic.gdx.Gdx
import me.grian.Main
import me.grian.network.NetworkManager
import me.grian.network.packets.c2s.C2SMovePacket

class Player(
    // TODO: maybe change this to a pos class for consistency with server, not sure tho cuz realx/y
    var x: Float,
    var y: Float,
    var realX: Float,
    var realY: Float,
    var name: String
) {
    fun move(x: Float, y: Float) {
        this.x = x.coerceIn(0.0f, (Gdx.graphics.width / Main.tileSize) - 1)
        this.y = y.coerceIn(0.0f, (Gdx.graphics.height / Main.tileSize) - 1)

        // maybe reconsider sending packet on function call, not sure though, have to see which one
        // ends up being more used, setting x/y directly or normnal move
        NetworkManager.sendPacket(
            C2SMovePacket(x.toInt(), y.toInt())
        )

        realX = this.x * Main.tileSize
        realY = this.y * Main.tileSize
    }

    constructor(x: Int, y: Int, name: String):
        this(x.toFloat(), y.toFloat(), x * Main.tileSize, y * Main.tileSize, name)
}
