package me.grian.packets.c2s

import me.grian.Clients
import me.grian.Constants
import me.grian.packets.PacketType

class C2SMovePacket : C2SPacket {
    override suspend fun handle(data: MutableMap<String, Any>, playerIdx: Int) {
        val x = data["x"]!! as Int
        val y = data["y"]!! as Int

        // should also be done client-side
        if (x !in 0..Constants.MAX_X || y !in 0..Constants.MAX_Y) return
        if (Clients.players.any { it.pos.x == x && it.pos.y == y }) return

        Clients.players[playerIdx].move(x, y)
        Clients.updateAllPlayers()
    }

    companion object {
        val STRUCTURE = mapOf(
            "x" to PacketType.INTEGER,
            "y" to PacketType.INTEGER
        )
    }
}