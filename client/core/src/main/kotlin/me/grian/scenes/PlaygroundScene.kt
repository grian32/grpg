package me.grian.scenes

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator.FreeTypeFontParameter
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.Main.Companion.player
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

            rect(player.realX, player.realY, tileSize, tileSize)
        }
    }

    private fun renderPlayerName(batch: SpriteBatch) {
        val nameYPos = if ((player.x == 0.0f || player.x == 1.0f) && player.y == 15.0f) {
            player.realY + (tileSize / 2)
        } else {
            player.realY + tileSize
        }

        blueFont.draw(batch, player.name, player.realX, nameYPos)
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
        redFont.draw(batch, "X: ${player.x} Y: ${player.y}", 0.0f, Gdx.graphics.height.toFloat())
    }
}
