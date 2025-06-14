package me.grian

import com.badlogic.gdx.ApplicationAdapter
import com.badlogic.gdx.Gdx
import com.badlogic.gdx.graphics.GL20
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import me.grian.player.Player
import me.grian.player.PlayerInputHandler
import me.grian.scenes.LoginScreenScene
import me.grian.scenes.PlaygroundScene

/** [com.badlogic.gdx.ApplicationListener] implementation shared by all platforms. */
class Main : ApplicationAdapter() {
    private lateinit var shapeRenderer: ShapeRenderer
    private lateinit var batch: SpriteBatch
    private lateinit var playerInputHandler: PlayerInputHandler

    override fun create() {
        shapeRenderer = ShapeRenderer()

        batch = SpriteBatch()

        playerInputHandler = PlayerInputHandler()

        PlaygroundScene.create()
        LoginScreenScene.create()
    }

    override fun render() {
        Gdx.gl.glClearColor(0.0f, 0.0f, 0.0f, 1.0f)
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT)


        if (isLoggedIn) {
            Gdx.input.inputProcessor = playerInputHandler
            PlaygroundScene.render(shapeRenderer, batch)
        } else {
            LoginScreenScene.render(shapeRenderer, batch)
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
