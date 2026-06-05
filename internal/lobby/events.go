package lobby

import (
	"game_bp/event_registry"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type LobbyInfo struct {
	EndTimer int64             `msg:"endTimer"`
	Players  []LobbyInfoPlayer `msg:"players"`
}

type LobbyInfoPlayer struct {
	Id    string `msg:"playerId"`
	Name  string `msg:"name"`
	Ready bool   `msg:"ready"`
}

type GameStart struct{}

var lobbyInfoSender = neoroute.Register[LobbyInfo](event_registry.EventReg, "lobby_info")
var gameStartSender = neoroute.Register[GameStart](event_registry.EventReg, "game_start")
