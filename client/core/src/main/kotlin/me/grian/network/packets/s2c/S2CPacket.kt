package me.grian.network.packets.s2c

import io.ktor.utils.io.*
import kotlinx.io.Buffer

interface S2CPacket {
    suspend fun handle(data: MutableMap<String, Any>)

    suspend fun decode(readChannel: ByteReadChannel): MutableMap<String, Any> {
        error("Decoding not supported for this type")
    }
}
