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
    private val clients: MutableMap<String, Player> = mutableMapOf()

    suspend fun handleClient(socket: Socket, receiveChannel: ByteReadChannel, sendChannel: ByteWriteChannel) {
        try {
            while (!receiveChannel.isClosedForRead) {
                val opcode = receiveChannel.readByte()

                val packet = c2sOpcodes.find { it.opcode == opcode }

                if (packet == null) {
                    logger.info("Received unknown opcode: $opcode from ${socket.localAddress}")
                    continue
                }

                // TODO: probably factor this out to a function or something it is getting a bit bulky i fear
                if (packet == C2SPacketOpcode.LOGIN) {
                    val strLength = receiveChannel.readInt()
                    val str = receiveChannel.readByteArray(strLength).toString(Charset.defaultCharset())
                    if (str !in clients) {
                        val allPositions = clients.values.map { it.pos }

                        var startingPoint = Point(0, 0)


                        // label bad :S & needs to be y first so it puts them on the bottom row first
                        outer@ for (y in 0..Constants.MAX_Y) {
                            for (x in 0..Constants.MAX_X) {
                                val pos = Point(x, y)
                                if (pos !in allPositions) {
                                    startingPoint = pos
                                    break@outer
                                }
                            }
                        }

                        val player = Player(
                            sendChannel,
                            startingPoint
                        )
                        clients[str] = player
                        // minX + miny on that row
                        sendToUser(str, S2CLoginAcceptedPacket(startingPoint))
                    } else {
                        sendToChannel(sendChannel, S2CLoginRejectedPacket())
                    }

                    logger.info("Client logged in with username [${str}]")
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
            val client = clients.entries.find { it.value.writeChannel == sendChannel }
            if (client == null) return
            clients.remove(client.key)
        }
    }

    suspend fun sendToUser(username: String, packet: S2CPacket) {
        val client = clients[username]?.writeChannel ?: return
        sendToChannel(client, packet)
    }

    suspend fun sendToChannel(sendChannel: ByteWriteChannel, packet: S2CPacket) {
        val buf = Buffer()

        buf.writeByte(packet.opcode)
        packet.handle(buf)
        sendChannel.writePacket(buf)
        sendChannel.flush()
    }
}