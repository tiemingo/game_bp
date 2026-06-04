package lobby

import (
	"game_bp/util/logger"
	"log/slog"
)

func (l *Lobby) SendPlayerInfoEvent() {

	players := []PlayerInfoPlayer{}
	for _, p := range l.players {
		players = append(players, PlayerInfoPlayer{
			Id:    p.id,
			Name:  p.name,
			Ready: p.ready,
		})
	}
	ev, err := playerInfoSender(PlayerInfo{
		Players: players,
	})
	if err != nil {
		slog.Info("failed to create player info event", logger.Err(err), logger.LobbyId(l.id))
		return
	}
	go func() {
		if err := l.adapterRegistry.Broadcast(ev); err != nil {
			slog.Info("failed to broadcast player info event", logger.Err(err), logger.LobbyId(l.id))
		}
	}()
}
