package game

import (
	"fmt"
	"time"
)

type GameState string

func (gs GameState) Stringer() string {
	return fmt.Sprintf("%s", gs)
}

func (gs GameState) String() string {
	return string(gs)
}

// GameState:
//
// 1. CREATED: Joindre une partie, en attente de tous les joueurs
// 2. SETUP: le maitre selectionne un mot, le insider recupere le mot
// 3. STARTED: Compte a rebours
// 4. STOPPED: Fin de partie, debat vote
// 5. RESOLVED: Revelation joueurs, role, mot
// 6. Option to reset to CREATED
//
const (
	CREATED  GameState = "CREATED"
	SETUP    GameState = "SETUP"
	STARTED  GameState = "STARTED"
	STOPPED  GameState = "STOPPED"
	RESOLVED GameState = "RESOLVED"
)

var phaseTransitions = make(map[GameState][]GameState)

func init() {
	// Initialize phase transitions contraints
	phaseTransitions[CREATED] = []GameState{RESOLVED}
	phaseTransitions[SETUP] = []GameState{CREATED}
	phaseTransitions[STARTED] = []GameState{SETUP}
	phaseTransitions[STOPPED] = []GameState{STARTED}
	phaseTransitions[RESOLVED] = []GameState{STOPPED}
}

type State struct {
	Name    string             `json:"name"`
	Phase   GameState          `json:"state"`
	Word    string             `json:"word"`
	Started time.Time          `json:"started_time"`
}

func NewState(name string) *State {
	return &State{
		Name:    name,
		Phase:   CREATED,
		Word:    GetNewWord(),
		Started: time.Time{},
	}
}

func (s *State) SetNewWord() string {
	s.Word = GetNewWord()
	return s.Word
}

func (s *State) SetStartTime() {
	s.Started = time.Now()
}

func (s *State) SetPhase(newPhase GameState) error {
	var changed = false
	for _, phase := range phaseTransitions[newPhase] {
		if phase == s.Phase {
			changed = true
			s.Phase = newPhase
		}
	}
	if !changed {
		return fmt.Errorf("Game phase not changed")
	} else {
		return nil
	}
}
