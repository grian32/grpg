package me.grian.network.packets.s2c

import me.grian.network.packets.PacketType
import kotlin.reflect.KClass

enum class S2CPacketOpcode(val opcode: Byte, val structure: Map<String, PacketType>, val packet: S2CPacket, val decodes: Boolean = false) {
    // basically a shim because i cant handle this normally lol
    LOGIN_ACCEPTED(0x01, S2CLoginAcceptedPacket.STRUCTURE, S2CLoginAcceptedPacket()),
    LOGIN_REJECTED(0x02, mapOf(), S2CLoginRejectedPacket()),
    PLAYERS_UPDATE(0x03, mapOf(), S2CPlayersUpdatePacket(), decodes = true)
}
