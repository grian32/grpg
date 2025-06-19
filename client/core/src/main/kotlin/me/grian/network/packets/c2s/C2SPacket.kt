package me.grian.network.packets.c2s

import kotlinx.io.Buffer

interface C2SPacket {
    val opcode: Byte

    fun handle(buf: Buffer)
}
