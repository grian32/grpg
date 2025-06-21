package me.grian.player

import io.ktor.utils.io.*
import kotlinx.io.Buffer
import me.grian.packets.s2c.S2CPacket
import me.grian.packets.s2c.sendToWriteChannel

data class Player(
    val name: String,
    val writeChannel: ByteWriteChannel,
    val pos: Point,
) {
    fun move(x: Int, y: Int) {
        pos.x = x
        pos.y = y
    }

    suspend fun sendPacket(packet: S2CPacket) {
        packet.sendToWriteChannel(writeChannel)
    }
}