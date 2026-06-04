package routes

import (
	"game_bp/handler/lobby_handler"
	"game_bp/internal/client"

	"github.com/Liphium/neoroute"
)

func SetupLobbyRoutes(r neoroute.Router[client.ClientData], t *neoroute.WebSocketTransporter[client.ClientData]) {
	handlerInfo := lobby_handler.HandlerInfo{
		T: t,
	}
	neoroute.Route(r, "create", handlerInfo.Create)
	neoroute.Route(r, "join", handlerInfo.Join)
	neoroute.RouteOk(r, "reconnect", handlerInfo.Reconnect)
	neoroute.RouteOk(r, "ready", lobby_handler.Ready)
	neoroute.RouteOkNoRequest(r, "leave", lobby_handler.Leave)
}
