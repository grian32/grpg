package me.grian.packets.c2s

interface C2SPacket {
    fun handle(data: MutableMap<String, Any>)
}