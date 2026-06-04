package lobby

import (
	"fmt"
	"game_bp/internal/client"
	"game_bp/util"

	"github.com/Liphium/neoroute"
	"github.com/tinylib/msgp/msgp"
)

func Modify(lobbyId string, modifyFunc func(*Lobby) error) error {
	l, ok := getLobby(lobbyId)
	if !ok || l.lobby == nil {
		return fmt.Errorf("%s", util.ErrLobbyNotFound)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return modifyFunc(l.lobby)
}

func ModifyNeo[RS any, PS interface {
	*RS
	msgp.Marshaler
}](c *neoroute.ResCtx[client.ClientData, RS, PS], lobbyId string, modifyFunc func(*Lobby) error) error {
	l, ok := getLobby(lobbyId)
	if !ok || l.lobby == nil {
		return c.RespondError(util.ErrLobbyNotFound)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return modifyFunc(l.lobby)
}

func ModifyNeoPlayer[RS any, PS interface {
	*RS
	msgp.Marshaler
}](c *neoroute.ResCtx[client.ClientData, RS, PS], modifyFunc func(*Lobby, *Player) error) error {

	// Get lobby ID from client data
	lobbyId := ""
	playerId := ""
	client.AccessData(&c.Ctx, func(cd *client.ClientData) error {
		lobbyId = cd.LobbyId
		playerId = cd.PlayerId
		return nil
	})

	return ModifyNeo(c, lobbyId, func(l *Lobby) error {

		// Check if player is part of the lobby
		p, exists := l.players[playerId]
		if !exists {
			return c.RespondError(util.ErrPlayerNotInLobby)
		}
		return modifyFunc(l, p)
	})
}

func ModifyNeoOk(c *neoroute.OkCtx[client.ClientData], lobbyId string, modifyFunc func(*Lobby) error) error {
	l, ok := getLobby(lobbyId)
	if !ok || l.lobby == nil {
		return c.RespondError(util.ErrLobbyNotFound)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return modifyFunc(l.lobby)
}

func ModifyNeoOkPlayer(c *neoroute.OkCtx[client.ClientData], modifyFunc func(*Lobby, *Player) error) error {

	// Get lobby ID from client data
	lobbyId := ""
	playerId := ""
	client.AccessData(&c.Ctx, func(cd *client.ClientData) error {
		lobbyId = cd.LobbyId
		playerId = cd.PlayerId
		return nil
	})

	return ModifyNeoOk(c, lobbyId, func(l *Lobby) error {

		// Check if player is part of the lobby
		p, exists := l.players[playerId]
		if !exists {
			return c.RespondError(util.ErrPlayerNotInLobby)
		}
		return modifyFunc(l, p)
	})
}
