package me.grian.player

import com.badlogic.gdx.Gdx
import me.grian.Main

class Player(
    var x: Float,
    var y: Float,
    var realX: Float,
    var realY: Float,
    var name: String
) {
    fun move(x: Float, y: Float) {
        this.x = x.coerceIn(0.0f, Gdx.graphics.width / Main.tileSize)
        this.y = y.coerceIn(0.0f, Gdx.graphics.height / Main.tileSize)
        realX = this.x * Main.tileSize
        realY = this.y * Main.tileSize
    }
}
