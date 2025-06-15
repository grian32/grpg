package me.grian.packets

import io.netty.buffer.ByteBuf

interface Packet {
    fun handle(data: ByteBuf)
}