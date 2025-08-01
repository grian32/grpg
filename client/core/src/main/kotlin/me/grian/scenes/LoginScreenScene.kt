package me.grian.scenes

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator.FreeTypeFontParameter
import com.badlogic.gdx.graphics.glutils.ShapeRenderer
import com.badlogic.gdx.scenes.scene2d.Actor
import com.badlogic.gdx.scenes.scene2d.InputEvent
import com.badlogic.gdx.scenes.scene2d.Stage
import com.badlogic.gdx.scenes.scene2d.ui.Table
import com.badlogic.gdx.scenes.scene2d.ui.TextButton
import com.badlogic.gdx.scenes.scene2d.ui.TextButton.TextButtonStyle
import com.badlogic.gdx.scenes.scene2d.ui.TextField
import com.badlogic.gdx.scenes.scene2d.ui.TextField.TextFieldStyle
import com.badlogic.gdx.scenes.scene2d.utils.ChangeListener
import com.badlogic.gdx.scenes.scene2d.utils.ClickListener
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import me.grian.Main
import me.grian.network.NetworkManager
import me.grian.network.packets.c2s.C2SLoginPacket
import me.grian.utils.createColorDrawable
import me.grian.utils.filledShape
import me.grian.utils.textButtonStyle
import me.grian.utils.textFieldStyle
import kotlin.math.log

object LoginScreenScene : Scene {
    private lateinit var font: BitmapFont
    private lateinit var titleFont: BitmapFont
    private lateinit var titleText: GlyphLayout
    private lateinit var enterNameText: GlyphLayout
    private lateinit var failedLoginText: GlyphLayout
    private lateinit var stage: Stage
    private lateinit var table: Table
    var shouldRenderFailedLoginText = false

    override fun create() {
        val generator = FreeTypeFontGenerator(Gdx.files.internal("ui/font.ttf"))
        val parameter = FreeTypeFontParameter()
        parameter.size = 24
        parameter.color = Color.WHITE
        font = generator.generateFont(parameter)

        parameter.size = 48

        titleFont = generator.generateFont(parameter)

        titleText = GlyphLayout(titleFont, "GRPG Client")
        enterNameText = GlyphLayout(font, "Enter Name Below:")
        failedLoginText = GlyphLayout(font, "Failed to login, most likely the name is already in use")

        stage = Stage()
        Gdx.input.inputProcessor = stage

        table = Table()
        table.setFillParent(true)
        stage.addActor(table)

        buildLayout()
    }

    override fun render(shapeRenderer: ShapeRenderer, batch: SpriteBatch) {
        renderTitleText(batch)
        renderLoginBox(shapeRenderer, batch)

        stage.act(Gdx.graphics.deltaTime)
        stage.draw()
    }

    override fun dispose() {
        font.dispose()
        titleFont.dispose()
        stage.dispose()
    }

    private fun renderTitleText(batch: SpriteBatch) {
        batch.begin()

        val textX = Gdx.graphics.width / 2.0f - titleText.width / 2.0f
        val textY = Gdx.graphics.height - 50.0f

        titleFont.draw(batch, titleText, textX, textY)

        batch.end()
    }

    private fun renderLoginBox(shapeRenderer: ShapeRenderer, batch: SpriteBatch) {
        val loginBgWidth = 400.0f
        val loginBgHeight = 200.0f

        val halfScreenWidth = Gdx.graphics.width / 2.0f
        val halfScreenHeight = Gdx.graphics.height / 2.0f

        filledShape(shapeRenderer) {
            color = Color.BROWN

            rect(
                halfScreenWidth - loginBgHeight,
                halfScreenHeight + loginBgHeight,
                loginBgWidth,
                loginBgHeight
            )
        }

        batch.begin()

        font.draw(
            batch,
            enterNameText,
            halfScreenHeight - enterNameText.width / 2.0f,
            halfScreenHeight + loginBgHeight * 2 - 25.0f
        )

        if (shouldRenderFailedLoginText) {
            font.draw(
                batch,
                failedLoginText,
                halfScreenWidth - failedLoginText.width / 2.0f,
                halfScreenHeight - loginBgHeight * 2
            )
        }

        batch.end()
    }

    private fun buildLayout() {
        val goldenrod = createColorDrawable(Color.GOLDENROD)
        val gold = createColorDrawable(Color.GOLD)

        val buttonStyle = textButtonStyle {
            font = this@LoginScreenScene.font
            up = goldenrod
            down = goldenrod
            over = gold
        }

        val textFieldStyle = textFieldStyle {
            font = this@LoginScreenScene.font
            fontColor = Color.WHITE

            background = goldenrod
            cursor = gold
            selection = gold
        }

        val field = TextField("", textFieldStyle)

        field.setTextFieldFilter { _, c -> c.isLetterOrDigit() }
        field.maxLength = 8

        val loginButton = TextButton("Login", buttonStyle)

        loginButton.addListener(object : ClickListener() {
            override fun clicked(event: InputEvent?, x: Float, y: Float) {
                Main.player.name = field.text
                NetworkManager.sendPacket(C2SLoginPacket(
                    field.text
                ))
            }
        })

        table.add(field).padBottom(45.0f).row()
        table.add(loginButton).padBottom(Gdx.graphics.height / 2.0f + 20.0f).row()
    }
}
