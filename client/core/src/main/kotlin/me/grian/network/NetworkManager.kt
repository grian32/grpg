package me.grian.network

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.io.Buffer
import me.grian.network.packets.PacketType
import me.grian.network.packets.c2s.C2SPacket
import me.grian.network.packets.s2c.S2CPacketOpcode
import org.slf4j.LoggerFactory
import java.nio.charset.Charset

object NetworkManager {
    private val selectorManager = SelectorManager(Dispatchers.IO)
    private lateinit var socket: Socket
    private val logger = LoggerFactory.getLogger(this::class.java)
    private lateinit var writeChannel: ByteWriteChannel

    private val s2cOpcodes = S2CPacketOpcode.entries

    val scope = CoroutineScope(Dispatchers.IO)

    suspend fun start() {
        // TODO don't hardcode these
        socket = aSocket(selectorManager).tcp().connect("127.0.0.1", 4422)

        val readChannel = socket.openReadChannel()
        writeChannel = socket.openWriteChannel(autoFlush = true)

        try {
            while (!readChannel.isClosedForRead) {
                val opcode = readChannel.readByte()

                val packet = s2cOpcodes.find { it.opcode == opcode }

                if (packet == null) {
                    logger.info("Received unknown opcode: $opcode from ${socket.localAddress}")
                    continue
                }

                packet.instance.handle(readChannel)
            }
        } catch (e: Throwable) {
            logger.error("Error reading from socket", e)
        }
    }

    fun dispose() {
        socket.close()
        selectorManager.close()
    }

    fun sendPacket(packet: C2SPacket) {
        scope.launch {
            val buf = Buffer()
            buf.writeByte(packet.opcode)
            packet.handle(buf)
            writeChannel.writePacket(buf)
            writeChannel.flush()
        }
    }
}
