package starter

import (
	"log/slog"
	"net/http"
	"os"

	"game_bp/routes"
	ws "game_bp/transporter/web_socket"
	"game_bp/util"
	"game_bp/util/logger"
)

func Start() {
	logger.InitLogger(os.Getenv("LOG_ENV"))

	slog.Info("Starting application", "version", util.VERSION)

	// Create WebSocket transporter
	wsHook, wsT := ws.CreateTransporter()
	wsT.AddEventRegistry(ws.EventReg)
	wsT.SetRouter(routes.SetupWSRoutes())

	// Create host transporter
	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHook)

	slog.Info("Starting server...", "listen", os.Getenv("LISTEN"))
	if err := http.ListenAndServe(os.Getenv("LISTEN"), mux); err != nil {
		slog.Error("Server failed", logger.Err(err))
		os.Exit(1)
	}
}
