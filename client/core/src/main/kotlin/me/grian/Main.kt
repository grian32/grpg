package me.grian

import com.badlogic.gdx.ApplicationAdapter
import com.badlogic.gdx.Gdx
import com.badlogic.gdx.Input
import com.badlogic.gdx.graphics.GL20
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.scenes.LoginScreenScene
import me.grian.scenes.PlaygroundScene
import kotlin.math.max
import kotlin.math.min

/** [com.badlogic.gdx.ApplicationListener] implementation shared by all platforms. */
class Main : ApplicationAdapter() {
    private lateinit var shapeRenderer: ShapeRenderer
    private lateinit var batch: SpriteBatch

    override fun create() {
        shapeRenderer = ShapeRenderer()

        batch = SpriteBatch()

        PlaygroundScene.create()
        LoginScreenScene.create()
    }

    override fun render() {
        Gdx.gl.glClearColor(0.0f, 0.0f, 0.0f, 1.0f)
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT)


        if (isLoggedIn) {
            handleInput()
            PlaygroundScene.render(shapeRenderer, batch)
        } else {
            LoginScreenScene.render(shapeRenderer, batch)
        }
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

    override fun dispose() {
        shapeRenderer.dispose()
        batch.dispose()
        PlaygroundScene.dispose()
        LoginScreenScene.dispose()

        for (i in texturesToDipose) {
            i.dispose()
        }
    }

    companion object {
        var tileSize = 64.0f
        val texturesToDipose = mutableListOf<Texture>()

        var isLoggedIn: Boolean = false

        var squareX = 0.0f
        var squareY = 0.0f
        var squarePosX = 0.0f
        var squarePosY = 0.0f
        var name = ""
    }
}
