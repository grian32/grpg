package me.grian.packets.c2s

import me.grian.packets.PacketType

class C2SLoginPacket : C2SPacket {
    override suspend fun handle(data: MutableMap<String, Any>, playerIdx: Int) {}

    companion object {
        val STRUCTURE = mapOf(
            "name" to PacketType.UTF8_STRING
        )
    }
}