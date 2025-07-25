package me.grian.scenes

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator.FreeTypeFontParameter
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.Main
import me.grian.Main.Companion.players
import me.grian.Main.Companion.tileSize
import me.grian.player.Player
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
        renderPlayer(shapeRenderer, Main.player)

        // TODO: use explicit
        for (i in players) {
            renderPlayer(shapeRenderer, i)
        }

        batch.begin()

        renderPlayerName(batch, Main.player)

        for (i in players) {
            renderPlayerName(batch, i)
        }

        renderCoordinates(batch)

        batch.end()
    }

    override fun dispose() {
        redFont.dispose()
        blueFont.dispose()
    }

    private fun renderPlayer(shapeRenderer: ShapeRenderer, player: Player) {
        filledShape(shapeRenderer) {
            color = Color.SKY

            rect(player.realX, player.realY, tileSize, tileSize)
        }
    }

    private fun renderPlayerName(batch: SpriteBatch, player: Player) {
        val nameYPos = if ((player.pos.x == 0 || player.pos.x == 1) && player.pos.y == 15) {
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
                // test code mainly so i can differentiate "chunks"
                if (Main.player.chunkPos.y == 1 && x == 8 && y == 8) {

                    filledShape(shapeRenderer) {
                        color = Color.BLACK
                        rect(x * tileSize, y * tileSize, tileSize, tileSize)
                    }
                    continue
                }

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
        redFont.draw(batch, "X: ${Main.player.pos.x} Y: ${Main.player.pos.y}", 0.0f, Gdx.graphics.height.toFloat())
    }
}
