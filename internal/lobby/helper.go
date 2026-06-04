package lobby

import (
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

func (l *Lobby) addPlayer(id, sessionId, name string) *Player {
	p := &Player{
		id:        id,
		sessionId: sessionId,
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

func isNameAllowed(name string) bool {
	return len(name) >= 3 && len(name) <= 10
}
