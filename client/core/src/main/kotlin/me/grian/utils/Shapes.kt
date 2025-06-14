package me.grian.utils

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.TextureRegion
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import com.badlogic.gdx.graphics.glutils.ShapeRenderer.ShapeType
import com.badlogic.gdx.scenes.scene2d.utils.Drawable
import com.badlogic.gdx.scenes.scene2d.utils.TextureRegionDrawable

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

fun createColorDrawable(color: Color): Drawable {
    val pixmap = Pixmap(1, 1, Pixmap.Format.RGBA8888)
    pixmap.setColor(color)
    pixmap.fill()
    val texture = Texture(pixmap)
    pixmap.dispose()
    return TextureRegionDrawable(TextureRegion(texture))
}
