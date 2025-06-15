package me.grian.packets

enum class PacketType(length: Int = 0) {
    UTF8_STRING,
    INTEGER(length = Int.SIZE_BYTES)
}