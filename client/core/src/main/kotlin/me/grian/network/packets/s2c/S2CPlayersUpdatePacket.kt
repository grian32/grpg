package me.grian.network.packets.s2c

import io.ktor.utils.io.*
import me.grian.Main
import me.grian.player.Player
import java.nio.charset.Charset

class S2CPlayersUpdatePacket : S2CPacket {
    override suspend fun handle(readChannel: ByteReadChannel) {
        val length = readChannel.readShort()

        val playerList = mutableListOf<Player>()

        repeat(length.toInt()) {
            val nameLen = readChannel.readInt()
            val name = readChannel.readByteArray(nameLen).toString(Charset.defaultCharset())
            val x = readChannel.readInt()
            val y = readChannel.readInt()

            if (name == Main.player.name) {
                // sending here cuz player move func also send packet, have to see
                // about that
                Main.player.pos.x = x
                Main.player.pos.y = y
                Main.player.realX = (x % Main.chunkSize) * Main.tileSize
                Main.player.realY = (y % Main.chunkSize) * Main.tileSize
                Main.player.chunkPos.x = x / Main.chunkSize
                Main.player.chunkPos.y = y / Main.chunkSize
            } else {
                playerList.add(Player(
                    x,
                    y,
                    x / Main.chunkSize,
                    y / Main.chunkSize,
                    name
                ))
            }

            Main.players = playerList
        }
    }
}
