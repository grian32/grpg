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
import me.grian.Main
import me.grian.utils.createColorDrawable
import me.grian.utils.filledShape
import kotlin.math.log

object LoginScreenScene : Scene {
    private lateinit var font: BitmapFont
    private lateinit var titleFont: BitmapFont
    private lateinit var titleText: GlyphLayout
    private lateinit var enterNameText: GlyphLayout
    private lateinit var stage: Stage
    private lateinit var table: Table

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

        stage = Stage()
        Gdx.input.inputProcessor = stage

        table = Table()
        table.setFillParent(true)
        stage.addActor(table)

        buildLayout()

        table.debug = true
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

        batch.end()
    }

    private fun buildLayout() {
        val goldenrod = createColorDrawable(Color.GOLDENROD)
        val gold = createColorDrawable(Color.GOLD)

        val buttonStyle = TextButtonStyle()

        buttonStyle.font = font
        buttonStyle.up = goldenrod
        buttonStyle.down = goldenrod
        buttonStyle.over = gold


        val textFieldStyle = TextFieldStyle()
        textFieldStyle.font = font
        textFieldStyle.fontColor = Color.WHITE

        textFieldStyle.background = goldenrod
        textFieldStyle.cursor = gold
        textFieldStyle.selection = gold

        val field = TextField("", textFieldStyle)

        field.setTextFieldFilter { _, c -> c.isLetterOrDigit() }
        field.maxLength = 8

        val loginButton = TextButton("Login", buttonStyle)

        loginButton.addListener(object : ClickListener() {
            override fun clicked(event: InputEvent?, x: Float, y: Float) {
                Main.name = field.text
                Main.isLoggedIn = true
            }
        })

        table.add(field).padBottom(45.0f).row()
        table.add(loginButton).padBottom(Gdx.graphics.height / 2.0f + 20.0f).row()
    }
}
