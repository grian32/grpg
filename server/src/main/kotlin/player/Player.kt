package me.grian.player

import io.ktor.utils.io.*

data class Player(
    val name: String,
    val writeChannel: ByteWriteChannel,
    val pos: Point,
)