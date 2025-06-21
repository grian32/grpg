package me.grian.packets.s2c

import io.ktor.utils.io.core.*
import kotlinx.io.Buffer
import me.grian.Clients

class S2CPlayersUpdatePacket : S2CPacket {
    override val opcode: Byte
        get() = 0x03

    override suspend fun handle(buf: Buffer) {
        val players = Clients.players

        buf.writeShort(players.size.toShort())

        for (i in players) {
            buf.writeInt(i.name.length)
            buf.writeFully(i.name.toByteArray())
            buf.writeInt(i.pos.x)
            buf.writeInt(i.pos.y)
        }
    }
}