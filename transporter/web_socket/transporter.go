package ws

import (
	"game_bp/internal/client"
	"game_bp/util/logger"
	"log/slog"
	"net/http"

	"github.com/Liphium/neoroute"
	"github.com/coder/websocket"
)

func CreateTransporter() (http.HandlerFunc, *neoroute.WebSocketTransporter[client.ClientData]) {
	hook, t := neoroute.NewWebSocketTransporter(neoroute.WSConfig[client.ClientData]{
		UpgradeFunc: websocket.Accept,
		OverwriteSessionFunc: func(id string) bool {
			return true
		},
		HandshakeFunc: func(r *http.Request) (*neoroute.Session[client.ClientData], bool) {
			return neoroute.NewSession[client.ClientData](client.NewSession()), true
		},
		EnterNetworkFunc: func(session *neoroute.Session[client.ClientData], t *neoroute.WebSocketTransporter[client.ClientData]) {

			slog.Info("user connected, creating adapter", logger.SessionId(session.Id()))

			// Add to adapter registry, in this case we don't have to manually unregister the adapter, because we want then in the registry until they disconnect.
			// Then they will be removed automatically.
			adapter, err := t.Adapt(session.Id())
			if err != nil {
				slog.Error("failed to create adapter", logger.SessionId(session.Id()), logger.Err(err))
				return
			}
			AdapterReg.Register(session.Id(), adapter)
		},
		DisconnectHandler: func(session *neoroute.Session[client.ClientData]) {
			slog.Info("client disconnected", logger.SessionId(session.Id()))
		},
	})
	return hook, t
}
