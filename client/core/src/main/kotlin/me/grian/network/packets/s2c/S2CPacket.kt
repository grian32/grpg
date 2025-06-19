package me.grian.network.packets.s2c

import kotlinx.io.Buffer

interface S2CPacket {
    suspend fun handle(data: MutableMap<String, Any>)
}
