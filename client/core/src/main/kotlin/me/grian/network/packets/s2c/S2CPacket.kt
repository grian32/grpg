package me.grian.network.packets.s2c

import io.ktor.utils.io.*
import kotlinx.io.Buffer

interface S2CPacket {
    suspend fun handle(readChannel: ByteReadChannel)
}
