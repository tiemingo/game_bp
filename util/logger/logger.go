package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// ctxKey is a private type to prevent context key collisions
type ctxKey string

const RequestIDKey ctxKey = "request_id"

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
			r.AddAttrs(slog.String("request_id", reqID))
		}
	}
	return h.Handler.Handle(ctx, r)
}

// InitLogger initializes the global slog logger.
// env: "production" for JSON structured logs, anything else for human-readable text logs.
func InitLogger(env string) {
	var baseHandler slog.Handler

	if env == "production" {
		// Production keeps clean, structured JSON
		baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		// Development gets beautiful ANSI colors and a clean layout
		baseHandler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.StampMilli,
			AddSource:  true,
		})
	}

	// Wrap it with our context interceptor (from the previous code)
	logger := slog.New(&ContextHandler{Handler: baseHandler})

	slog.SetDefault(logger)
}

// Err is a clean helper to format errors consistently in your logs
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any("error", err.Error())
}

// SessionId is a clean helper to format session IDs consistently in your logs
func SessionId(sessionId string) slog.Attr {
	return slog.String("session_id", sessionId)
}

// LobbyId is a clean helper to format lobby IDs consistently in your logs
func LobbyId(lobbyId string) slog.Attr {
	return slog.String("lobby_id", lobbyId)
}

// Phase is a clean helper to format phase names consistently in your logs
func Phase(phase string) slog.Attr {
	return slog.String("phase", phase)
}
