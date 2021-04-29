package session

import (
	"insider/game"
	"time"
)

type GameState struct {
	Name      string                  `json:"name"`
	Phase     string                  `json:"phase"`
	Word      string                  `json:"word"`
	Started   time.Time               `json:"started_time"`
	CurrentId string                  `json:"current_id"`
	Players   []*PlayerState          `json:"players"`
}

type PlayerState struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Presenter struct {
	gs  GameState
	ctx SessionContext
}

func (s SessionContext) BuildPresenter() Presenter {
	gs := GameState{
		Name:      s.Session.Game.Name,
		Phase:     s.Session.Game.Phase.String(),
		Word:      "",
		Started:   s.Session.Game.Started,
		CurrentId: s.ID.String(),
	}

	for id, p := range s.Session.Players {
		playerState := PlayerState{
			Id:   "",
			Name: p.Name,
			Role: "",
		}
		// Allow extra information on the current player
		// requesting informations
		if id == s.ID {
			playerState.Id = p.Id.String()
			playerState.Role = p.Role.String()
		}
		// Allow all players to know who is the game master
		if p.Role == game.MASTER {
			playerState.Role = p.Role.String()
		}

		gs.Players = append(gs.Players, &playerState)
	}

	return Presenter{
		gs:  gs,
		ctx: s,
	}
}

func (p *Presenter) WithPlayersRoles() *Presenter {
	for id, player := range p.ctx.Session.Players {
		playerState := PlayerState{
			Id:   "",
			Name: player.Name,
			Role: "",
		}
		// Allow extra information on the current player
		// requesting informations
		if id == p.ctx.ID {
			playerState.Id = player.Id.String()
		}
		// All players can see all roles
		playerState.Role = player.Role.String()

		p.gs.Players = append(p.gs.Players, &playerState)
	}
	return p
}

func (p *Presenter) WithWord() *Presenter {
	p.gs.Word = p.ctx.Session.Game.Word
	return p
}

func (p Presenter) Build() GameState {
	return p.gs
}
