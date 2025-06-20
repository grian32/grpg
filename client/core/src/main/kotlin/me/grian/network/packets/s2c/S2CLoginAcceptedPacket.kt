package me.grian.network.packets.s2c

import me.grian.Main
import me.grian.network.packets.PacketType
import me.grian.scenes.LoginScreenScene

class S2CLoginAcceptedPacket : S2CPacket {
    override suspend fun handle(data: MutableMap<String, Any>) {
        Main.isLoggedIn = true
        LoginScreenScene.shouldRenderFailedLoginText = false

        val initialX = data["initialX"]!! as Int
        val initialY = data["initialY"]!! as Int

        Main.player.move(initialX.toFloat(), initialY.toFloat())
    }

    companion object {
        val STRUCTURE = mapOf(
            "initialX" to PacketType.INTEGER,
            "initialY" to PacketType.INTEGER
        )
    }
}
