package lobby_handler

import (
	"game_bp/internal/client"

	"github.com/Liphium/neoroute"
)

type HandlerInfo struct {
	T *neoroute.WebSocketTransporter[client.ClientData]
}
