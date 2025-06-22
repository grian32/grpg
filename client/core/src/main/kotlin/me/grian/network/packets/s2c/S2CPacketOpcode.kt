package me.grian.network.packets.s2c

import me.grian.network.packets.PacketType

enum class S2CPacketOpcode(val opcode: Byte, val instance: S2CPacket) {
    // basically a shim because i cant handle this normally lol
    LOGIN_ACCEPTED(0x01, S2CLoginAcceptedPacket()),
    LOGIN_REJECTED(0x02, S2CLoginRejectedPacket()),
    PLAYERS_UPDATE(0x03, S2CPlayersUpdatePacket())
}
