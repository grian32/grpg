package me.grian.network.packets.c2s

import io.ktor.utils.io.core.*
import kotlinx.io.Buffer
import me.grian.network.packets.c2s.C2SPacket

class C2SLoginPacket(
    val playerName: String
) : C2SPacket {
    override val opcode: Byte
        get() = 0x01

    override fun handle(buf: Buffer) {
        buf.writeInt(playerName.length)
        buf.writeFully(playerName.toByteArray())
    }
}
