package me.grian.utils

import com.badlogic.gdx.scenes.scene2d.ui.TextButton
import com.badlogic.gdx.scenes.scene2d.ui.TextButton.TextButtonStyle
import com.badlogic.gdx.scenes.scene2d.ui.TextField.TextFieldStyle

fun textButtonStyle(func: TextButtonStyle.() -> Unit): TextButtonStyle {
    val style = TextButtonStyle()

    style.func()

    return style
}

fun textFieldStyle(func: TextFieldStyle.() -> Unit): TextFieldStyle {
    val style = TextFieldStyle()

    style.func()

    return style
}
