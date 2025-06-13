package me.grian

import com.badlogic.gdx.ApplicationAdapter
import com.badlogic.gdx.Gdx
import com.badlogic.gdx.Input
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.GL20
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator.FreeTypeFontParameter
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.utils.filledShape
import me.grian.utils.lineShape
import kotlin.math.max
import kotlin.math.min

/** [com.badlogic.gdx.ApplicationListener] implementation shared by all platforms. */
class Main : ApplicationAdapter() {
    private lateinit var shapeRenderer: ShapeRenderer
    private lateinit var redFont: BitmapFont
    private lateinit var blueFont: BitmapFont
    private lateinit var batch: SpriteBatch


    private var tileSize = 64.0f
    private var squareX = 0.0f
    private var squareY = 0.0f
    private var squarePosX = 0.0f
    private var squarePosY = 0.0f
    private var name = "Grian"

    override fun create() {
        shapeRenderer = ShapeRenderer()

        val generator = FreeTypeFontGenerator(Gdx.files.internal("ui/font.ttf"))
        val parameter = FreeTypeFontParameter()
        parameter.size = 24
        parameter.color = Color.RED
        redFont = generator.generateFont(parameter)

        parameter.color = Color.BLUE
        blueFont = generator.generateFont(parameter)

        batch = SpriteBatch()
    }

    override fun render() {
        Gdx.gl.glClearColor(0.0f, 0.0f, 0.0f, 1.0f)
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT)

        handleInput()

        renderGrid()

        renderPlayer()

        batch.begin()

        renderPlayerName()

        renderCoordinates()

        batch.end()

    }

    private fun handleInput() {
        when {
            Gdx.input.isKeyJustPressed(Input.Keys.W) -> squareY += tileSize
            Gdx.input.isKeyJustPressed(Input.Keys.S) -> squareY -= tileSize
            Gdx.input.isKeyJustPressed(Input.Keys.A) -> squareX -= tileSize
            Gdx.input.isKeyJustPressed(Input.Keys.D) -> squareX += tileSize
        }

        squareY = max(0.0f, min(Gdx.graphics.height - tileSize, squareY))
        squarePosY = squareY / tileSize
        squareX = max(0.0f, min(Gdx.graphics.width - tileSize, squareX))
        squarePosX = squareX / tileSize
    }

    private fun renderPlayer() {
        filledShape(shapeRenderer) {
            color = Color.SKY

            rect(squareX, squareY, tileSize, tileSize)
        }
    }

    private fun renderPlayerName() {
        val nameYPos = if ((squarePosX == 0.0f || squarePosX == 1.0f) && squarePosY == 15.0f) {
            squareY + (tileSize / 2)
        } else {
            squareY + tileSize
        }

        blueFont.draw(batch, name, squareX, nameYPos)
    }

    private fun renderGrid() {
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

    private fun renderCoordinates() {
        redFont.draw(batch, "X: ${squareX / tileSize} Y: ${squareY / tileSize}", 0.0f, Gdx.graphics.height.toFloat())
    }

    override fun dispose() {
        shapeRenderer.dispose()
        redFont.dispose()
        blueFont.dispose()
        batch.dispose()
    }
}
