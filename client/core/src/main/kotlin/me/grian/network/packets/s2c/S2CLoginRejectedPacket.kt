package me.grian.network.packets.s2c

import io.ktor.utils.io.*
import me.grian.Main
import me.grian.scenes.LoginScreenScene

class S2CLoginRejectedPacket : S2CPacket {
    override suspend fun handle(readChannel: ByteReadChannel) {
        // more safety than anything lol
        Main.isLoggedIn = false
        LoginScreenScene.shouldRenderFailedLoginText = true
    }
}
