package me.grian.packets.s2c

import io.ktor.utils.io.*
import kotlinx.io.Buffer

suspend fun S2CPacket.sendToWriteChannel(writeChannel: ByteWriteChannel) {
    val buf = Buffer()

    buf.writeByte(opcode)
    handle(buf)
    writeChannel.writePacket(buf)
    writeChannel.flush()
}