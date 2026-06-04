package lobby

import (
	"sync"

	"github.com/google/uuid"
)

func getLobby(lobbyId string) (*wrappedLobby, bool) {
	value, ok := lobbies.Load(lobbyId)
	if !ok {
		return &wrappedLobby{}, false
	}
	return value.(*wrappedLobby), true
}

func reserveLobbyId() string {
	for {
		id := uuid.NewString()
		_, loaded := lobbies.LoadOrStore(id, &wrappedLobby{lobby: nil})
		if !loaded {
			return id
		}
	}
}

func addLobby(lobbyId string, lobby *Lobby) {
	lobbies.Store(lobbyId, &wrappedLobby{lobby: lobby, mutex: &sync.Mutex{}})
}

func deleteLobby(lobbyId string) {
	lobbies.Delete(lobbyId)
}
