package session

import (
	"context"
	"net/http"

	"insider/game"
	. "insider/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const Key string = "Session"

type SessionContext struct {
	Session    Session
	ID         uuid.UUID
	NewSession bool
}

func Inject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		gameName := vars["name"]
		if len(gameName) < 3 {
			JSON(w, Struct{"err": "Game name too short"}, http.StatusBadRequest)
			return
		}

		session, new := sessions.GetOrCreateSession(gameName)

		id, err := uuid.Parse(vars["uuid"])
		if err != nil {
			id = uuid.Nil
		}

		sessionContext := SessionContext{
			Session:    session,
			ID:         id,
			NewSession: new,
		}
		req := r.WithContext(context.WithValue(r.Context(), Key, sessionContext))

		next.ServeHTTP(w, req)
	})
}

func (s SessionContext) IsIdentified() bool {
	if s.ID == uuid.Nil {
		return false
	} else {
		_, ok := s.Session.Players.PlayerById(s.ID)
		return ok
	}
}

func (s SessionContext) IsMaster() bool {
	if !s.IsIdentified() {
		return false
	}

	return s.Session.Players.IsRole(s.ID, game.MASTER)
}

func (s SessionContext) IsInsider() bool {
	if !s.IsIdentified() {
		return false
	}

	return s.Session.Players.IsRole(s.ID, game.INSIDER)
}
