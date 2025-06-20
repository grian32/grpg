package me.grian

import me.grian.player.Point

fun findFirstAvailablePosition(allPositions: Set<Point>): Point {
    for (y in 0..Constants.MAX_Y) {
        for (x in 0..Constants.MAX_X) {
            val pos = Point(x, y)
            if (pos !in allPositions) {
                return pos
            }
        }
    }

    // fallback
    return Point(0, 0)
}