package me.grian.packets.s2c

import kotlinx.io.Buffer
import me.grian.player.Point

class S2CLoginAcceptedPacket(
    private val initialPos: Point
) : S2CPacket {
    override val opcode: Byte
        get() = 0x01

    override suspend fun handle(buf: Buffer) {
        buf.writeInt(initialPos.x)
        buf.writeInt(initialPos.y)
    }
}