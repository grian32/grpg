package me.grian.player

import com.badlogic.gdx.Input
import com.badlogic.gdx.InputProcessor
import me.grian.Main.Companion.player

class PlayerInputHandler : InputProcessor {
    override fun keyDown(keycode: Int): Boolean {
        return when (keycode) {
            Input.Keys.W -> { player.move(player.x, player.y + 1); true }
            Input.Keys.S -> { player.move(player.x, player.y - 1); true }
            Input.Keys.A -> { player.move(player.x - 1, player.y); true }
            Input.Keys.D -> { player.move(player.x + 1, player.y); true }
            else -> false
        }
    }

    override fun keyUp(keycode: Int): Boolean = false
    override fun keyTyped(character: Char): Boolean = false
    override fun touchDown(screenX: Int, screenY: Int, pointer: Int, button: Int): Boolean = false
    override fun touchUp(screenX: Int, screenY: Int, pointer: Int, button: Int): Boolean = false
    override fun touchCancelled(screenX: Int, screenY: Int, pointer: Int, button: Int): Boolean = false
    override fun touchDragged(screenX: Int, screenY: Int, pointer: Int): Boolean = false
    override fun mouseMoved(screenX: Int, screenY: Int): Boolean = false
    override fun scrolled(amountX: Float, amountY: Float): Boolean = false
}
