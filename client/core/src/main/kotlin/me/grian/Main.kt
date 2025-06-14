package me.grian

import com.badlogic.gdx.ApplicationAdapter
import com.badlogic.gdx.Gdx
import com.badlogic.gdx.Input
import com.badlogic.gdx.graphics.GL20
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.player.Player
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
            Gdx.input.isKeyJustPressed(Input.Keys.W) -> player.move(player.x, player.y + 1)
            Gdx.input.isKeyJustPressed(Input.Keys.S) -> player.move(player.x, player.y - 1)
            Gdx.input.isKeyJustPressed(Input.Keys.A) -> player.move(player.x - 1, player.y)
            Gdx.input.isKeyJustPressed(Input.Keys.D) -> player.move(player.x + 1, player.y)
        }
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

        val player = Player(0.0f, 0.0f, 0.0f, 0.0f, "")
    }
}
