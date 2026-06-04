package lobby_handler

import (
	"game_bp/internal/client"
	"game_bp/internal/lobby"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type ReadyRequest struct {
	Ready bool `msg:"ready"`
}

func Ready(c *neoroute.OkCtx[client.ClientData], req ReadyRequest) error {
	return lobby.ModifyNeoOkPlayer(c, func(l *lobby.Lobby, p *lobby.Player) error {

		err := l.Ready(p, req.Ready)
		if err != "" {
			return c.RespondError(err)
		}

		// Send player info event
		l.SendPlayerInfoEvent()

		return c.RespondOk()
	})
}
