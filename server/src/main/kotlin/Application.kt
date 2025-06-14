package me.grian

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import org.slf4j.LoggerFactory
import kotlin.math.log


fun main(args: Array<String>) {
    val logger = LoggerFactory.getLogger("main")

    runBlocking {
        val selectorManager = SelectorManager(Dispatchers.IO)
        val serverSocket = aSocket(selectorManager).tcp().bind("127.0.0.1", 4422)
        logger.info("Server is listening at ${serverSocket.localAddress}")
        while (true) {
            val socket = serverSocket.accept()
            logger.info("Accepted $socket")

            launch {
                val receiveChannel = socket.openReadChannel()
                val sendChannel = socket.openWriteChannel(autoFlush = true)

                try {
                    while (!receiveChannel.isClosedForRead) {
                        logger.info(receiveChannel.readUTF8Line() ?: break)
                    }
                } catch (e: Throwable) {
                    logger.error("Error reading from socket", e)
                } finally {
                    socket.close()
                }
            }
        }
    }
}