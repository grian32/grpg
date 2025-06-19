package me.grian.network.packets.s2c

import me.grian.Main
import me.grian.scenes.LoginScreenScene

class S2CLoginAcceptedPacket : S2CPacket {
    override suspend fun handle(data: MutableMap<String, Any>) {
        Main.isLoggedIn = true
        LoginScreenScene.shouldRenderFailedLoginText = false
    }
}
