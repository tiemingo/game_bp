package lobby

import (
	"game_bp/internal/phase"
	"game_bp/util"

	"github.com/google/uuid"
)

func (l *Lobby) newPlayerId() string {
	for {
		id := uuid.NewString()
		if _, exists := l.players[id]; !exists {
			return id
		}
	}
}

func (l *Lobby) addPlayer(id, name string) *Player {
	p := &Player{
		id:        id,
		name:      name,
		connected: true,
		token:     uuid.NewString(),
	}
	l.players[id] = p
	return p
}

func (l *Lobby) removePlayer(p *Player) {
	delete(l.players, p.id)
}

func (l *Lobby) isNameUnique(name string) bool {
	for _, p := range l.players {
		if p.name == name {
			return false
		}
	}
	return true
}

// canChangeReady checks if the ready status of a player can be changed.
// If it is allowed an empty string is returned, otherwise an error message is returned.
func (l *Lobby) canChangeReady() string {
	if !l.isRunning {
		return ""
	}
	return util.ErrLobbyRunning
}

// shouldSkipPhase checks if phase should be ended prematurely.
func (l *Lobby) shouldSkipPhase(ts phase.TimerStatus) bool {
	if ts.CurrentPhase == PHASE_LOBBY {
		return l.allPlayersReady() && len(l.players) >= l.minPlayers
	}

	// TODO: add other conditions for skipping phases if needed
	return false
}

func isNameAllowed(name string) bool {
	return len(name) >= 3 && len(name) <= 10
}

func (l *Lobby) allPlayersReady() bool {
	for _, p := range l.players {
		if !p.ready {
			return false
		}
	}
	return true
}

func (l *Lobby) stopOrStartLobbyTimer(add bool) {

	if len(l.players) == l.minPlayers && add {

		// Start lobby timer
		doneChan := make(chan struct{})
		l.commandChan <- phase.Command{
			Type: phase.CmdResumeIf,
			ResumeIf: func(ts phase.TimerStatus) bool {
				return ts.CurrentPhase == PHASE_LOBBY
			},
			DoneChan: doneChan,
		}
		<-doneChan
		return
	}

	if len(l.players) < l.minPlayers && !add {
		// Stop lobby timer
		doneChan := make(chan struct{})
		l.commandChan <- phase.Command{
			Type: phase.CmdResetIf,
			ResetIf: func(ts phase.TimerStatus) bool {
				return ts.CurrentPhase == PHASE_LOBBY
			},
			DoneChan: doneChan,
		}
		<-doneChan
		return
	}
}
