package me.grian.packets.s2c

import io.ktor.utils.io.*

interface S2CPacket {
    val opcode: Byte

    suspend fun handle(sendChannel: ByteWriteChannel)
}
