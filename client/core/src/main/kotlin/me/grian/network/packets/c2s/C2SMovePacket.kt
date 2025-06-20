package me.grian.network.packets.c2s

import kotlinx.io.Buffer

class C2SMovePacket(
    val x: Int,
    val y: Int
) : C2SPacket {
    override val opcode: Byte
        get() = 0x02

    override fun handle(buf: Buffer) {
        buf.writeInt(x)
        buf.writeInt(y)
    }
}
