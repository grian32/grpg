package me.grian.scenes

import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.glutils.ShapeRenderer

interface Scene {
    fun create()
    fun render(shapeRenderer: ShapeRenderer, batch: SpriteBatch)
    fun dispose()
}
