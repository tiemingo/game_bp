package lobby_handler

import (
	"game_bp/internal/client"
	"game_bp/internal/lobby"
	"game_bp/util"

	"github.com/Liphium/neoroute"
)

func Leave(c *neoroute.OkCtx[client.ClientData]) error {

	return client.AccessData(&c.Ctx, func(cd *client.ClientData) error {

		// Check if player is already in a lobby
		if cd.LobbyId == "" {
			return c.RespondError(util.ErrPlayerNotInLobby)
		}

		return lobby.ModifyNeoOkPlayer(c, func(l *lobby.Lobby, p *lobby.Player) error {

			err := l.Leave(p, c.Session().Id())
			if err != "" {
				return c.RespondError(err)
			}

			// Clear client data
			cd.LobbyId = ""
			cd.PlayerId = ""

			// Send player info event
			l.SendPlayerInfoEvent()

			return c.RespondOk()
		})
	})

}
