package me.grian

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import io.netty.buffer.Unpooled
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import me.grian.packets.c2s.C2SPacket
import me.grian.packets.PacketType
import me.grian.packets.c2s.C2SPacketOpcode
import org.slf4j.LoggerFactory
import kotlin.reflect.full.primaryConstructor

fun main(args: Array<String>) {
    val logger = LoggerFactory.getLogger("main")

    runBlocking {
        val selectorManager = SelectorManager(Dispatchers.IO)
        val serverSocket = aSocket(selectorManager).tcp().bind("127.0.0.1", 4422)
        logger.info("Server is listening at ${serverSocket.localAddress}")

        while (true) {
            val socket = serverSocket.accept()
            logger.info("Accepted ${socket.localAddress}")

            launch {
                val receiveChannel = socket.openReadChannel()
                val sendChannel = socket.openWriteChannel(autoFlush = true)

                Clients.handleClient(socket, receiveChannel, sendChannel)
            }
        }
    }
}