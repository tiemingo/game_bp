package lobby_handler

import (
	"game_bp/internal/client"
	"game_bp/internal/lobby"
	"game_bp/util"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type CreateRequest struct {
	Name string `msg:"name"`
}

type CreateResponse struct {
	LobbyId     string `msg:"lobbyId"`
	LobbyToken  string `msg:"lobbyToken"`
	PlayerId    string `msg:"playerId"`
	PlayerToken string `msg:"playerToken"`
}

func (h HandlerInfo) Create(c *neoroute.ResCtx[client.ClientData, CreateResponse, *CreateResponse], req CreateRequest) error {

	return client.AccessData(&c.Ctx, func(cd *client.ClientData) error {

		// Check if player is already in a lobby
		if cd.LobbyId != "" {
			return c.RespondError(util.ErrPlayerAlreadyInLobby)
		}

		adapter, adaptErr := h.GetAdapterFunc(c.Session().Id())
		if adaptErr != nil {
			return c.RespondError(util.ErrInternalServerError)
		}

		// Create lobby
		lobbyId, lobbyToken, playerId, playerToken, err := lobby.CreateLobby(c.Session().Id(), adapter, req.Name)
		if err != "" {
			return c.RespondError(err)
		}

		// Set client data
		cd.LobbyId = lobbyId
		cd.PlayerId = playerId

		return c.Respond(CreateResponse{
			LobbyId:     lobbyId,
			LobbyToken:  lobbyToken,
			PlayerId:    playerId,
			PlayerToken: playerToken,
		})
	})
}
