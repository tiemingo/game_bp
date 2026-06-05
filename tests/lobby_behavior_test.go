package tests

import (
	"fmt"
	"game_bp/event_registry"
	"game_bp/handler/lobby_handler"
	"game_bp/internal/client"
	"game_bp/util"
	"testing"

	"github.com/Liphium/neoroute"
	"github.com/stretchr/testify/assert"
)

func TestLobbyBehavior(t *testing.T) {

	lobbyId := ""
	lobbyToken := ""
	neo := neoroute.NewNeoRouter[client.ClientData](neoroute.Config{})

	type TestPlayer[D any] struct {
		SessionId string
		Session   *neoroute.Session[D]
		Id        string
		Token     string
		Adapter   neoroute.Adapter
		Name      string
	}

	creator := TestPlayer[client.ClientData]{
		SessionId: "creator-session-id",
		Adapter:   neoroute.NewTestingAdapter([]*neoroute.EventRegistry{event_registry.EventReg}),
		Name:      "Creator",
	}
	creator.Session = neoroute.NewTestingSession[client.ClientData](client.ClientData{}, creator.SessionId)

	joiningPlayer := TestPlayer[client.ClientData]{
		SessionId: "joining-player-session-id",
		Adapter:   neoroute.NewTestingAdapter([]*neoroute.EventRegistry{event_registry.EventReg}),
		Name:      "Joining",
	}
	joiningPlayer.Session = neoroute.NewTestingSession[client.ClientData](client.ClientData{}, joiningPlayer.SessionId)

	handlerInfo := lobby_handler.HandlerInfo{
		GetAdapterFunc: func(sessionId string) (neoroute.Adapter, error) {
			if sessionId == creator.SessionId {
				return creator.Adapter, nil
			} else if sessionId == joiningPlayer.SessionId {
				return joiningPlayer.Adapter, nil
			}
			return nil, fmt.Errorf("adapter not found for session ID: %s", sessionId)
		},
	}

	t.Run("create lobby valid", func(t *testing.T) {
		ctx := neoroute.NewTestingResCtx[client.ClientData, lobby_handler.CreateResponse](neo, "lobby.create", creator.Session)

		_, errMsg, err := neoroute.GetTestingResponse[lobby_handler.CreateResponse](handlerInfo.Create(ctx, lobby_handler.CreateRequest{
			Name: "Name is too long. It should be less than 10 characters.",
		}))
		assert.Nil(t, err)
		assert.Equal(t, util.ErrInvalidName, errMsg)
	})

	t.Run("create lobby with invalid name", func(t *testing.T) {
		ctx := neoroute.NewTestingResCtx[client.ClientData, lobby_handler.CreateResponse](neo, "lobby.create", creator.Session)

		resp, errMsg, err := neoroute.GetTestingResponse[lobby_handler.CreateResponse](handlerInfo.Create(ctx, lobby_handler.CreateRequest{
			Name: creator.Name,
		}))
		assert.Nil(t, err)
		assert.Empty(t, errMsg)
		lobbyId = resp.LobbyId
		lobbyToken = resp.LobbyToken
		creator.Id = resp.PlayerId
		creator.Token = resp.PlayerToken
	})

	t.Run("join lobby", func(t *testing.T) {
		ctx := neoroute.NewTestingResCtx[client.ClientData, lobby_handler.JoinResponse](neo, "lobby.join", joiningPlayer.Session)

		resp, errMsg, err := neoroute.GetTestingResponse[lobby_handler.JoinResponse](handlerInfo.Join(ctx, lobby_handler.JoinRequest{
			Name:       joiningPlayer.Name,
			LobbyId:    lobbyId,
			LobbyToken: lobbyToken,
		}))
		assert.Nil(t, err)
		assert.Empty(t, errMsg)
		joiningPlayer.Id = resp.PlayerId
		joiningPlayer.Token = resp.PlayerToken
	})
	_ = lobbyId
	_ = lobbyToken
}
