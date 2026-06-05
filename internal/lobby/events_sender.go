package lobby

import (
	"game_bp/internal/phase"
	"game_bp/util/logger"
	"log/slog"
)

func (l *Lobby) SendLobbyInfoEvent() {

	// Collect player info
	players := []LobbyInfoPlayer{}
	for _, p := range l.players {
		players = append(players, LobbyInfoPlayer{
			Id:    p.id,
			Name:  p.name,
			Ready: p.ready,
		})
	}

	// Get end timestamp
	replyChan := make(chan phase.TimerStatus)
	l.commandChan <- phase.Command{
		Type:      phase.CmdGetTimerStatus,
		ReplyChan: replyChan,
	}
	timerStatus := <-replyChan

	ev, err := lobbyInfoSender(LobbyInfo{
		EndTimer: timerStatus.EndTime.Unix(),
		Players:  players,
	})
	if err != nil {
		slog.Info("failed to create lobby info event", logger.Err(err), logger.LobbyId(l.id))
		return
	}
	go func() {
		if err := l.adapterRegistry.Broadcast(ev); err != nil {
			slog.Info("failed to broadcast lobby info event", logger.Err(err), logger.LobbyId(l.id))
		}
	}()
}

func (l *Lobby) SendGameStartEvent() {

	ev, err := gameStartSender(GameStart{})
	if err != nil {
		slog.Info("failed to create game start event", logger.Err(err), logger.LobbyId(l.id))
		return
	}
	go func() {
		if err := l.adapterRegistry.Broadcast(ev); err != nil {
			slog.Info("failed to broadcast game start event", logger.Err(err), logger.LobbyId(l.id))
		}
	}()
}
