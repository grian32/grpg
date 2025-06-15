package me.grian.packets.c2s

import io.netty.buffer.ByteBuf
import me.grian.packets.Packet
import java.nio.charset.Charset

class C2SLoginPacket : Packet {
    override fun handle(data: ByteBuf) {
        println("here")
        val strLength = data.readInt()
        val str = data.readBytes(strLength)

        println(str.toString(Charset.defaultCharset()))
    }
}