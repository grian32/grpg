package me.grian.packets.c2s

interface C2SPacket {
    /**
     * @param data A map containing the packet data defined in C2SPacketOpcode
     * @param playerIdx the index of the player who sent this request in the player list in Clients
     */
    fun handle(data: MutableMap<String, Any>, playerIdx: Int)
}