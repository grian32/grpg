package me.grian

import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import kotlinx.io.Buffer
import me.grian.packets.PacketType
import me.grian.packets.c2s.C2SPacket
import me.grian.packets.c2s.C2SPacketOpcode
import me.grian.packets.s2c.S2CLoginAcceptedPacket
import me.grian.packets.s2c.S2CLoginRejectedPacket
import me.grian.packets.s2c.S2CPacket
import me.grian.player.Player
import me.grian.player.Point
import org.slf4j.LoggerFactory
import java.nio.charset.Charset
import kotlin.reflect.full.primaryConstructor

object Clients {
    private val c2sOpcodes = C2SPacketOpcode.entries
    private val logger = LoggerFactory.getLogger(this::class.java)
    private val clients: MutableList<Player> = mutableListOf()

    suspend fun handleClient(socket: Socket, receiveChannel: ByteReadChannel, writeChannel: ByteWriteChannel) {
        try {
            while (!receiveChannel.isClosedForRead) {
                val opcode = receiveChannel.readByte()

                val packet = c2sOpcodes.find { it.opcode == opcode }

                if (packet == null) {
                    logger.info("Received unknown opcode: $opcode from ${socket.localAddress}")
                    continue
                }

                if (packet == C2SPacketOpcode.LOGIN) {
                    processLogin(receiveChannel, writeChannel)
                    println(clients)
                    continue
                }

                val packetData = mutableMapOf<String, Any>()

                for ((name, dataType) in packet.structure) {
                    val data = when (dataType) {
                        PacketType.UTF8_STRING -> {
                            val strLength = receiveChannel.readInt()
                            val str = receiveChannel.readByteArray(strLength).toString(Charset.defaultCharset())

                            str
                        }

                        PacketType.INTEGER -> receiveChannel.readInt()
                        PacketType.BYTE -> receiveChannel.readByte()
                    }

                    packetData[name] = data
                }

                val instance = packet.packet.primaryConstructor!!.call() as C2SPacket
                instance.handle(packetData)
            }
        } catch (e: Throwable) {
            logger.error("Error reading from socket", e)
            val client = clients.find { it.writeChannel == writeChannel }
            if (client == null) return
            clients.remove(client)
        }
    }

    suspend fun processLogin(receiveChannel: ByteReadChannel, writeChannel: ByteWriteChannel) {
        val strLength = receiveChannel.readInt()
        val str = receiveChannel.readByteArray(strLength).toString(Charset.defaultCharset())
        val names = clients.map { it.name }
        if (str !in names) {
            var startingPoint = Point(0, 0)

            if (clients.isNotEmpty()) {
                startingPoint = findFirstAvailablePosition(
                    clients.mapTo(mutableSetOf()) { it.pos }
                )
            }

            val player = Player(
                str,
                writeChannel,
                startingPoint
            )
            clients.add(player)
            sendToUser(str, S2CLoginAcceptedPacket(startingPoint))
        } else {
            sendToChannel(writeChannel, S2CLoginRejectedPacket())
        }

        logger.info("Client logged in with username [${str}]")
        println(clients)
    }

    suspend fun sendToUser(username: String, packet: S2CPacket) {
        val client = clients.find { it.name == username } ?: return
        sendToChannel(client.writeChannel, packet)
    }

    suspend fun sendToChannel(writeChannel: ByteWriteChannel, packet: S2CPacket) {
        val buf = Buffer()

        buf.writeByte(packet.opcode)
        packet.handle(buf)
        writeChannel.writePacket(buf)
        writeChannel.flush()
    }
}