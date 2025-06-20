package me.grian.player

import io.ktor.utils.io.*

data class Player(
    val name: String,
    val writeChannel: ByteWriteChannel,
    val pos: Point,
) {
    fun move(x: Int, y: Int) {
        pos.x = x
        pos.y = y
    }
}