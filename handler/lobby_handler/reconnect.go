package lobby_handler

import (
	"game_bp/internal/client"
	"game_bp/internal/lobby"
	"game_bp/util"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type ReconnectRequest struct {
	LobbyId     string `msg:"lobbyId"`
	PlayerId    string `msg:"playerId"`
	PlayerToken string `msg:"playerToken"`
}

func Reconnect(c *neoroute.OkCtx[client.ClientData], req ReconnectRequest) error {

	return client.AccessData(&c.Ctx, func(cd *client.ClientData) error {

		// Check if player is already in a lobby
		if cd.LobbyId != "" {
			return c.RespondError(util.ErrPlayerAlreadyInLobby)
		}
		return lobby.Modify(req.LobbyId, func(l *lobby.Lobby) error {

			err := l.Reconnect(req.PlayerId, c.Session().Id(), req.PlayerToken)
			if err != "" {
				return c.RespondError(err)
			}

			// Set client data
			cd.LobbyId = req.LobbyId
			cd.PlayerId = req.PlayerId

			return c.RespondOk()
		})

	})
}
