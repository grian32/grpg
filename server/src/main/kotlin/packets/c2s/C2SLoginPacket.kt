package me.grian.packets.c2s

import io.netty.buffer.ByteBuf
import java.nio.charset.Charset

class C2SLoginPacket : C2SPacket {
    override fun handle(data: ByteBuf) {
        println("here")
        val strLength = data.readInt()
        val str = data.readBytes(strLength)

        println(str.toString(Charset.defaultCharset()))
    }
}