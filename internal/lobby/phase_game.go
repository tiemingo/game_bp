package lobby

import (
	"game_bp/internal/phase"
	"game_bp/util/logger"
	"log/slog"
	"time"
)

func (p *phaseInfo) phaseGameEnd() (phase.Phase, time.Duration, bool) {
	slog.Debug("phase ended", logger.Phase("game"), logger.LobbyId(p.lobbyId))

	Modify(p.lobbyId, func(l *Lobby) error {
		slog.Debug("phase end lobby update", logger.Phase("game"), logger.LobbyId(p.lobbyId))

		l.round++
		return nil
	})

	slog.Debug("phase end finished continuing to next phase",
		logger.Phase("game"),
		"next_phase", PHASE_GAME,
		logger.LobbyId(p.lobbyId),
	)

	return PHASE_GAME, PHASE_GAME_DURATION, true
}
