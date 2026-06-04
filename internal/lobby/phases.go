package lobby

import (
	"game_bp/internal/phase"
	"time"
)

const (
	PHASE_LOBBY phase.Phase = "PHASE_LOBBY"
	PHASE_GAME  phase.Phase = "PHASE_GAME"
)

const (
	PHASE_LOBBY_DURATION = 5 * time.Minute
	PHASE_GAME_DURATION  = 5 * time.Second
)

type phaseInfo struct {
	lobbyId string
}
