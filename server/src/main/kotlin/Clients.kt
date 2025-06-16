package me.grian

import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import me.grian.packets.c2s.C2SPacketOpcode
import org.slf4j.LoggerFactory
import java.nio.charset.Charset

object Clients {
    private val c2sOpcodes = C2SPacketOpcode.entries
    private val logger = LoggerFactory.getLogger(this::class.java)
    private val clients: MutableMap<String, ByteWriteChannel> = mutableMapOf()

    suspend fun handleClient(socket: Socket, receiveChannel: ByteReadChannel, sendChannel: ByteWriteChannel) {
        try {
            while (!receiveChannel.isClosedForRead) {
                println(clients)

                val opcode = receiveChannel.readByte()

                val packet = c2sOpcodes.find { it.opcode == opcode }

                if (packet == null) {
                    logger.info("Received unknown opcode: $opcode from ${socket.localAddress}")
                    continue
                }

                if (packet == C2SPacketOpcode.LOGIN) {
                    val strLength = receiveChannel.readInt()
                    val str = receiveChannel.readByteArray(strLength).toString(Charset.defaultCharset())
                    clients[str] = sendChannel
                    continue
                }

            }
        } catch (e: Throwable) {
            logger.error("Error reading from socket", e)
        } finally {
            socket.close()
        }
    }
}