package util

var (
	ErrLobbyNotFound          = "Lobby was not found."
	ErrPlayerNotInLobby       = "Player is not in lobby."
	ErrLobbyRunning           = "Lobby is already running."
	ErrNameTaken              = "Player with this name already exists."
	ErrInvalidName            = "Name must be between 3 and 10 characters long."
	ErrInvalidPlayerToken     = "Invalid player token."
	ErrPlayerAlreadyConnected = "Player is already connected."
	ErrPlayerKicked           = "Player was kicked from the lobby."
	ErrPlayerAlreadyInLobby   = "Player is already in a lobby."
	ErrReadyStatusUnchanged   = "Ready status is already set to the given value."
)
