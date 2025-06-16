package me.grian.packets.s2c

import kotlinx.io.Buffer

interface S2CPacket {
    val opcode: Byte

    suspend fun handle(buf: Buffer)
}
