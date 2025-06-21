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
import me.grian.packets.s2c.sendToWriteChannel
import me.grian.player.Player
import me.grian.player.Point
import org.slf4j.LoggerFactory
import java.nio.charset.Charset
import kotlin.reflect.full.primaryConstructor

object Clients {
    private val c2sOpcodes = C2SPacketOpcode.entries
    private val logger = LoggerFactory.getLogger(this::class.java)
    val players: MutableList<Player> = mutableListOf()

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
                    continue
                }

                val client = players.indexOfFirst { it.writeChannel == writeChannel }
                if (client == -1) return

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


                packet.packet.handle(packetData, client)
                println(players)
            }
        } catch (e: Throwable) {
            logger.error("Error reading from socket", e)
            val client = players.find { it.writeChannel == writeChannel }
            if (client == null) return
            players.remove(client)
        }
    }

    suspend fun processLogin(receiveChannel: ByteReadChannel, writeChannel: ByteWriteChannel) {
        val strLength = receiveChannel.readInt()
        val str = receiveChannel.readByteArray(strLength).toString(Charset.defaultCharset())
        val names = players.map { it.name }
        if (str !in names) {
            var startingPoint = Point(0, 0)

            if (players.isNotEmpty()) {
                startingPoint = findFirstAvailablePosition(
                    players.mapTo(mutableSetOf()) { it.pos }
                )
            }

            val player = Player(
                str,
                writeChannel,
                startingPoint
            )
            players.add(player)
            player.sendPacket(S2CLoginAcceptedPacket(startingPoint))
        } else {
            S2CLoginRejectedPacket().sendToWriteChannel(writeChannel)
        }

        logger.info("Client logged in with username [${str}]")
        println(players)
    }
}