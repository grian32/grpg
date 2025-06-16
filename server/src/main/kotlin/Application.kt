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

        val c2sOpcodes = C2SPacketOpcode.entries

        while (true) {
            val socket = serverSocket.accept()
            logger.info("Accepted ${socket.localAddress}")

            launch {
                val receiveChannel = socket.openReadChannel()
                val sendChannel = socket.openWriteChannel(autoFlush = true)

                try {
                    while (!receiveChannel.isClosedForRead) {
                        val opcode = receiveChannel.readByte()

                        val packet = c2sOpcodes.find { it.opcode == opcode }

                        if (packet == null) {
                            logger.info("Received unknown opcode: $opcode from ${socket.localAddress}")
                        } else {
                            if (packet.structure.contains(PacketType.UTF8_STRING)) {
                                val bytes = Unpooled.buffer()
                                for (i in packet.structure) {
                                    if (i == PacketType.UTF8_STRING) {
                                        val len = receiveChannel.readInt()
                                        val byteArray = ByteArray(len)
                                        receiveChannel.readFully(byteArray)

                                        bytes.writeInt(len)
                                        bytes.writeBytes(byteArray)
                                    } else {
                                        // TODO OTHER HANDLING
                                    }
                                }

                                val inst = packet.packet.primaryConstructor!!.call() as C2SPacket
                                inst.handle(bytes)
                                continue
                            } else {
                                // TODO: send packet lengths etc
                            }
                        }
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