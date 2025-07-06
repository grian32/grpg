package shared

import "net"

type Player struct {
	X uint32
	Y uint32
	// might not need these will see how design pans out
	ChunkX uint32
	ChunkY uint32
	Name   string
	Conn   net.Conn
}
