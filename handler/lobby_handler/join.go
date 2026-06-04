package lobby_handler

import (
	"game_bp/internal/client"
	"game_bp/internal/lobby"
	"game_bp/util"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type JoinRequest struct {
	Name       string `msg:"name"`
	LobbyId    string `msg:"lobbyId"`
	LobbyToken string `msg:"lobbyToken"`
}

type JoinResponse struct {
	PlayerId string `msg:"playerId"`
	Token    string `msg:"tokenToken"`
}

func (h HandlerInfo) Join(c *neoroute.ResCtx[client.ClientData, JoinResponse, *JoinResponse], req JoinRequest) error {
	return client.AccessData(&c.Ctx, func(cd *client.ClientData) error {

		// Check if player is already in a lobby
		if cd.LobbyId != "" {
			return c.RespondError("Player is already in a lobby.")
		}

		return lobby.Modify(req.LobbyId, func(l *lobby.Lobby) error {

			adapter, adaptErr := h.T.Adapt(c.Session().Id())
			if adaptErr != nil {
				return c.RespondError(util.ErrInternalServerError)
			}

			// Join lobby
			playerId, playerToken, err := l.Join(c.Session().Id(), adapter, req.Name)
			if err != "" {
				return c.RespondError(err)
			}

			// Set client data
			cd.LobbyId = req.LobbyId
			cd.PlayerId = playerId

			// Send player info event
			l.SendPlayerInfoEvent()

			return c.Respond(JoinResponse{
				PlayerId: playerId,
				Token:    playerToken,
			})
		})
	})
}
