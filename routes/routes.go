package routes

import (
	"game_bp/internal/client"
	"game_bp/util/logger"
	"log/slog"

	"github.com/Liphium/neoroute"
)

func SetupWSRoutes() *neoroute.NeoRouter[client.ClientData] {

	// Create router
	r := neoroute.NewNeoRouter[client.ClientData](neoroute.Config{
		ErrorHandler: func(err error) string {
			slog.Info("error on WebSocket router", logger.Err(err))
			return "Internal server error."
		},
	})

	SetupLobbyRoutes(r.Group("lobby"))

	return r
}
