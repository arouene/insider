package game

import (
	"math/rand"
)

var (
	words []string
)

func InitWords(w []string) {
	words = w
}

func GetNewWord() string {
	listSize := len(words)
	return words[rand.Intn(listSize)]
}
