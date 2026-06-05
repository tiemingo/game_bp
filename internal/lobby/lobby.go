package lobby

import (
	"game_bp/internal/phase"
	"game_bp/util"

	"github.com/Liphium/neoroute"
)

type Lobby struct {
	id              string                    // Constant field
	token           string                    // Constant field
	phaseManager    *phase.PhaseManager       // Constant field
	doneChan        chan struct{}             // Constant field
	commandChan     chan<- phase.Command      // Constant field
	adapterRegistry *neoroute.AdapterRegistry // Constant field
	maxPlayers      int                       // Constant field
	minPlayers      int                       // Constant field

	// Mutable fields
	isRunning bool
	round     int

	players map[string]*Player
}

func (l *Lobby) Join(sessionId string, adapter neoroute.Adapter, name string) (string, string, string) {
	if l.isRunning {
		return "", "", util.ErrLobbyRunning
	}

	if !l.isNameUnique(name) {
		return "", "", util.ErrNameTaken
	}

	if len(l.players) >= l.maxPlayers {
		return "", "", util.ErrLobbyFull
	}

	p := l.addPlayer(l.newPlayerId(), name)

	// Add player sessionId to adapter registry
	l.adapterRegistry.Register(sessionId, adapter)

	return p.id, p.token, ""
}

func (l *Lobby) Ready(p *Player, ready bool) string {

	if err := l.canChangeReady(); err != "" {
		return err
	}

	if p.ready == ready {
		return util.ErrReadyStatusUnchanged
	}

	p.ready = ready

	// Skip phase if wanted
	skipDoneChan := make(chan struct{})
	l.commandChan <- phase.Command{
		Type:     phase.CmdSkipIf,
		DoneChan: skipDoneChan,
		SkipIf:   l.shouldSkipPhase,
	}
	<-skipDoneChan

	return ""
}

func (l *Lobby) Reconnect(id string, sessionId string, adapter neoroute.Adapter, playerToken string) string {

	p, ok := l.players[id]
	if !ok {
		return util.ErrPlayerNotInLobby
	}

	if p.token != playerToken {
		return util.ErrInvalidPlayerToken
	}

	if p.connected {
		return util.ErrPlayerAlreadyConnected
	}

	if p.kicked {
		return util.ErrPlayerKicked
	}

	p.connected = true

	// Add player sessionId to adapter registry
	l.adapterRegistry.Register(sessionId, adapter)

	return ""
}

func (l *Lobby) Leave(p *Player, sessionId string) string {

	if l.isRunning {
		return util.ErrLobbyRunning
	}

	l.removePlayer(p)

	// Remove player sessionId from adapter registry
	l.adapterRegistry.Unregister(sessionId)

	return ""
}
