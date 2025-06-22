package me.grian.network.packets.s2c

import io.ktor.utils.io.*

class S2CPlayersUpdatePacket : S2CPacket {
    override suspend fun handle(data: MutableMap<String, Any>) {
        TODO("Not yet implemented")
    }

    override suspend fun decode(readChannel: ByteReadChannel): MutableMap<String, Any> {
        TODO("Not yet implemented")
    }
}
