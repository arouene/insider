package game

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type PlayerRole string

func (p PlayerRole) Stringer() string {
	return fmt.Sprintf("%s", p)
}

func (p PlayerRole) String() string {
	return string(p)
}

const (
	MASTER  PlayerRole = "MASTER"
	INSIDER PlayerRole = "INSIDER"
	CITIZEN PlayerRole = "CITIZEN"
	NONE    PlayerRole = "NONE"
)

// Player is a player of the game
// it has a name a unique uuid and a role
type Player struct {
	Id   uuid.UUID  `json:"id"`
	Name string     `json:"name"`
	Role PlayerRole `json:"role"`
}

func NewPlayer(name string, id uuid.UUID) Player {
	return Player{
		Id:   id,
		Name: name,
		Role: NONE,
	}
}

func (p *Player) SetRole(role PlayerRole) *Player {
	p.Role = role
	return p
}

func (p Player) IsInsider() bool {
	return p.Role == INSIDER
}

func (p Player) IsMaster() bool {
	return p.Role == MASTER
}

// Players is a map of Player
// it contains all the players of a game
// each player is identified by a unique uuid
type Players map[uuid.UUID]*Player

func NewsPlayers() Players {
	return make(Players)
}

// AddPlayer add a player in the players list
// will not fail, but update player if Id already exists
func (p Players) AddPlayer(player Player) error {
	p[player.Id] = &player

	return nil
}

func (p Players) PlayerById(id uuid.UUID) (*Player, bool) {
	if player, ok := p[id]; ok {
		return player, true
	}
	return &Player{}, false
}

func (p Players) PlayerByName(name string) (*Player, bool) {
	for _, player := range p {
		if player.Name == name {
			return player, true
		}
	}
	return &Player{}, false
}

func (p Players) SetRandomRoles() error {
	if len(p) < 4 {
		return fmt.Errorf("Not enough player, there must be at least 4 players")
	}

	master := rand.Intn(len(p))
	insider := master
	for insider == master {
		insider = rand.Intn(len(p))
	}

	var i = 0
	for _, v := range p {
		switch i {
		case master:
			v.SetRole(MASTER)
		case insider:
			v.SetRole(INSIDER)
		default:
			v.SetRole(CITIZEN)
		}
		i++
	}

	return nil
}

func (p Players) IsRole(id uuid.UUID, r PlayerRole) bool {
	if p, ok := p.PlayerById(id); ok {
		return p.Role == r
	} else {
		return false
	}
}
