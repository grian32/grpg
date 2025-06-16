package me.grian.packets.c2s

import me.grian.packets.PacketType
import kotlin.reflect.KClass

enum class C2SPacketOpcode(val opcode: Byte, val structure: Map<String, PacketType>, val packet: KClass<*>) {
    // basically a shim because i cant handle this normally lol
    LOGIN(0x01, C2SLoginPacket.STRUCTURE, C2SLoginPacket::class)
}