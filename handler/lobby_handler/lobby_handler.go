package lobby_handler

import (
	"github.com/Liphium/neoroute"
)

type HandlerInfo struct {
	GetAdapterFunc func(sessionId string) (neoroute.Adapter, error)
}
