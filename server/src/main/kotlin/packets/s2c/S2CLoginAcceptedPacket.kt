package me.grian.packets.s2c

import kotlinx.io.Buffer

class S2CLoginAcceptedPacket : S2CPacket {
    override val opcode: Byte
        get() = 0x01

    override suspend fun handle(buf: Buffer) {
        // empty packet
    }
}