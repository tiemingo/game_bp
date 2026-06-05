package lobby

import (
	"game_bp/event_registry"

	"github.com/Liphium/neoroute"
)

//go:generate msgp

type PlayerInfo struct {
	Players []PlayerInfoPlayer `msg:"players"`
}

type PlayerInfoPlayer struct {
	Id    string `msg:"playerId"`
	Name  string `msg:"name"`
	Ready bool   `msg:"ready"`
}

type GameStart struct{}

var playerInfoSender = neoroute.Register[PlayerInfo](event_registry.EventReg, "player_info")
var gameStartSender = neoroute.Register[GameStart](event_registry.EventReg, "game_start")
