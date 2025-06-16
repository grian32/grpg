package me.grian.packets.s2c

import io.ktor.utils.io.*

class S2CLoginAcceptedPacket : S2CPacket {
    override val opcode: Byte
        get() = 0x01

    override suspend fun handle(sendChannel: ByteWriteChannel) {
        // empty packet
    }
}