package tests

import (
	"fmt"
	"game_bp/event_registry"
	"game_bp/handler/lobby_handler"
	"game_bp/internal/client"
	"game_bp/internal/lobby"
	"game_bp/util"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/stretchr/testify/assert"
)

func TestLobbyBehavior(t *testing.T) {

	t.Run("lobby behavior test", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {

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
					switch sessionId {
					case creator.SessionId:
						return creator.Adapter, nil
					case joiningPlayer.SessionId:
						return joiningPlayer.Adapter, nil
					}
					return nil, fmt.Errorf("adapter not found for session ID: %s", sessionId)
				},
			}

			// Create valid lobby
			func() {
				ctx := neoroute.NewTestingResCtx[client.ClientData, lobby_handler.CreateResponse](neo, "lobby.create", creator.Session)

				_, errMsg, err := neoroute.GetTestingResponse[lobby_handler.CreateResponse](handlerInfo.Create(ctx, lobby_handler.CreateRequest{
					Name: "Name is too long. It should be less than 10 characters.",
				}))
				assert.Nil(t, err)
				assert.Equal(t, util.ErrInvalidName, errMsg)
			}()

			// create lobby with invalid name
			func() {
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
			}()

			// join lobby
			func() {
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

				// Verify that events are sent to both players
				time.Sleep(time.Millisecond)
				events, err := creator.Adapter.(*neoroute.TestingAdapter).DrainEvents()
				assert.Nil(t, err)
				assert.Len(t, events, 1)
				assert.Equal(t, "lobby_info", events[0].Name)
				ev, err := neoroute.UnmarshalEventTesting[lobby.LobbyInfo](events[0].Data)
				assert.Nil(t, err)
				assert.Len(t, ev.Players, 2)
				assert.Contains(t, []string{ev.Players[0].Name, ev.Players[1].Name}, creator.Name)
				assert.Contains(t, []string{ev.Players[0].Name, ev.Players[1].Name}, joiningPlayer.Name)
			}()

			// Stop the lobby
			func() {
				assert.Nil(t, lobby.Modify(lobbyId, func(l *lobby.Lobby) error {
					l.Stop()
					return nil
				}))
			}()
		})
	})
}
