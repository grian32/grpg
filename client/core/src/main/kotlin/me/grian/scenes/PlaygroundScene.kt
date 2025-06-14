package me.grian.scenes

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator.FreeTypeFontParameter
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.Main.Companion.name
import me.grian.Main.Companion.squarePosX
import me.grian.Main.Companion.squarePosY
import me.grian.Main.Companion.squareX
import me.grian.Main.Companion.squareY
import me.grian.Main.Companion.tileSize
import me.grian.utils.filledShape
import me.grian.utils.lineShape

object PlaygroundScene : Scene {
    private lateinit var redFont: BitmapFont
    private lateinit var blueFont: BitmapFont

    override fun create() {
        val generator = FreeTypeFontGenerator(Gdx.files.internal("ui/font.ttf"))
        val parameter = FreeTypeFontParameter()
        parameter.size = 24
        parameter.color = Color.RED
        redFont = generator.generateFont(parameter)

        parameter.size = 16
        parameter.color = Color.BLUE
        blueFont = generator.generateFont(parameter)
    }

    override fun render(shapeRenderer: ShapeRenderer, batch: SpriteBatch) {
        renderGrid(shapeRenderer)
        renderPlayer(shapeRenderer)

        batch.begin()

        renderPlayerName(batch)
        renderCoordinates(batch)

        batch.end()
    }

    override fun dispose() {
        redFont.dispose()
        blueFont.dispose()
    }

    private fun renderPlayer(shapeRenderer: ShapeRenderer) {
        filledShape(shapeRenderer) {
            color = Color.SKY

            rect(squareX, squareY, tileSize, tileSize)
        }
    }

    private fun renderPlayerName(batch: SpriteBatch) {
        val nameYPos = if ((squarePosX == 0.0f || squarePosX == 1.0f) && squarePosY == 15.0f) {
            squareY + (tileSize / 2)
        } else {
            squareY + tileSize
        }

        blueFont.draw(batch, name, squareX, nameYPos)
    }

    private fun renderGrid(shapeRenderer: ShapeRenderer) {
        val gridX = (Gdx.graphics.width / tileSize).toInt()
        val gridY = (Gdx.graphics.height / tileSize).toInt()

        for (x in 0..gridX) {
            for (y in 0..gridY) {
                filledShape(shapeRenderer) {
                    color = Color.WHITE

                    rect(x * tileSize, y * tileSize, tileSize, tileSize)
                }

                lineShape(shapeRenderer) {
                    color = Color.BLACK

                    rect(x * tileSize, y * tileSize, 4.0f, 4.0f)
                }
            }
        }
    }

    private fun renderCoordinates(batch: SpriteBatch) {
        redFont.draw(batch, "X: ${squareX / tileSize} Y: ${squareY / tileSize}", 0.0f, Gdx.graphics.height.toFloat())
    }
}
