package routes

import (
	"game_bp/handler/lobby_handler"
	"game_bp/internal/client"

	"github.com/Liphium/neoroute"
)

func SetupLobbyRoutes(r neoroute.Router[client.ClientData]) {
	neoroute.Route(r, "create", lobby_handler.Create)
	neoroute.Route(r, "join", lobby_handler.Join)
	neoroute.RouteOk(r, "reconnect", lobby_handler.Reconnect)
	neoroute.RouteOk(r, "ready", lobby_handler.Ready)
	neoroute.RouteOkNoRequest(r, "leave", lobby_handler.Leave)
}
