package me.grian.packets.c2s

import me.grian.packets.PacketType
import kotlin.reflect.KClass

enum class C2SPacketOpcode(val opcode: Byte, val structure: List<PacketType>, val packet: KClass<*>) {
    LOGIN(0x01, listOf(PacketType.UTF8_STRING), C2SLoginPacket::class)
}