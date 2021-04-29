package session

import (
	"time"

	"insider/game"
)

var (
	sessions Sessions
)

func init() {
	sessions = NewSessions()
}

type Session struct {
	Game    *game.State
	Players game.Players
	Started time.Time
}

func NewSession(name string) Session {
	return Session{
		Game:    game.NewState(name),
		Players: game.NewsPlayers(),
		Started: time.Now(),
	}
}

// Sessions is a map of multiple sessions
type Sessions map[string]Session

func NewSessions() Sessions {
	return make(Sessions)
}

// GetOrCreateSession returns true if the session is a new one
func (s Sessions) GetOrCreateSession(name string) (Session, bool) {
	if _, ok := s[name]; ok {
		return s[name], false
	}

	s[name] = NewSession(name)

	return s[name], true
}
