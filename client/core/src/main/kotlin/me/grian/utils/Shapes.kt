package me.grian.utils

import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import com.badlogic.gdx.graphics.glutils.ShapeRenderer.ShapeType

fun filledShape(sr: ShapeRenderer, func: ShapeRenderer.() -> Unit) {
    sr.begin(ShapeType.Filled)

    sr.func()

    sr.end()
}

fun lineShape(sr: ShapeRenderer, func: ShapeRenderer.() -> Unit) {
    sr.begin(ShapeType.Line)

    sr.func()

    sr.end()
}
