package me.grian.network.packets.s2c

import io.ktor.utils.io.*
import me.grian.Main
import me.grian.scenes.LoginScreenScene

class S2CLoginAcceptedPacket : S2CPacket {
    override suspend fun handle(readChannel: ByteReadChannel) {
        Main.isLoggedIn = true
        LoginScreenScene.shouldRenderFailedLoginText = false

        val initialX = readChannel.readInt()
        val initialY = readChannel.readInt()

        Main.player.move(initialX.toFloat(), initialY.toFloat())
    }
}
