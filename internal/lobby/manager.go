package lobby

import (
	"game_bp/internal/phase"
	"game_bp/util"
	"sync"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/google/uuid"
)

type wrappedLobby struct {
	lobby *Lobby
	mutex *sync.Mutex
}

var lobbies = &sync.Map{} // map[string]*wrappedLobby

func CreateLobby(sessionId string, adapter neoroute.Adapter, name string) (string, string, string, string, string) {

	if !isNameAllowed(name) {
		return "", "", "", "", util.ErrInvalidName
	}

	id := reserveLobbyId()

	// Initialize phase manager
	pi := phaseInfo{lobbyId: id}
	phaseManager, doneChan, commandChan, err := phase.NewPhaseManager(phase.Config{
		InitialPhase:    PHASE_LOBBY,
		InitialDuration: PHASE_LOBBY_DURATION,
		Phases: map[phase.Phase]func() (phase.Phase, time.Duration, bool){
			PHASE_LOBBY: pi.phaseLobbyEnd,
			PHASE_GAME:  pi.phaseGameEnd,
		},
	})
	if err != nil {
		panic(err)
	}

	// Create lobby
	l := &Lobby{
		id:              id,
		token:           uuid.NewString(),
		phaseManager:    phaseManager,
		commandChan:     commandChan,
		doneChan:        doneChan,
		adapterRegistry: neoroute.NewAdapterRegistry(),
		maxPlayers:      4, // TODO: change max players
		minPlayers:      2, // TODO: change min players

		players: make(map[string]*Player),
	}

	// Add initial player
	p := l.addPlayer(l.newPlayerId(), name)

	// Add player sessionId to adapter registry
	l.adapterRegistry.Register(sessionId, adapter)

	// Start phase manager
	go phaseManager.RunEngine()

	// Stop and reset the timer
	resetChan := make(chan struct{})
	l.commandChan <- phase.Command{
		Type: phase.CmdResetIf,
		ResetIf: func(status phase.TimerStatus) bool {
			return true
		},
		DoneChan: resetChan,
	}
	<-resetChan

	addLobby(id, l)
	return id, l.token, p.id, p.token, ""
}
