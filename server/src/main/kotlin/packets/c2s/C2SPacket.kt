package me.grian.packets.c2s

import io.netty.buffer.ByteBuf

interface C2SPacket {
    fun handle(data: ByteBuf)
}