package client

import (
	"sync"

	"github.com/google/uuid"
)

var (
	sessions   = map[string]struct{}{}
	sessionsMu sync.RWMutex
)

func AddSession(sessionId string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	sessions[sessionId] = struct{}{}
}

func RemoveSession(sessionId string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	delete(sessions, sessionId)
}

func SessionExists(sessionId string) bool {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()

	_, exists := sessions[sessionId]
	return exists
}

func NewSession() string {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	for {
		id := uuid.NewString()
		if _, exists := sessions[id]; !exists {
			sessions[id] = struct{}{}
			return id
		}
	}
}
