package lobby

type Player struct {
	id    string // Constant field
	name  string // Constant field
	token string // Constant field

	sessionId string
	connected bool
	kicked    bool
	ready     bool
}
