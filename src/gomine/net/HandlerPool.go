package net

import (
	"gomine/players"
	"gomine/net/info"
	"gomine/players/handlers"
)

var registeredHandlers = map[int][]players.IPacketHandler{}

func InitHandlerPool() {
	RegisterPacketHandler(info.LoginPacket, handlers.NewLoginHandler())
}

func RegisterPacketHandler(id int, handler players.IPacketHandler) {
	registeredHandlers[id] = append(registeredHandlers[id], handler)
}

func GetPacketHandlers(id int) []players.IPacketHandler {
	return registeredHandlers[id]
}