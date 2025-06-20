package me.grian.player

import io.ktor.utils.io.*

data class Player(
    val writeChannel: ByteWriteChannel,
    val pos: Point,
)