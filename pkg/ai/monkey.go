package ai

import (
	"math/rand"

	"github.com/ola-rozenfeld/kulami/pkg/board"
)

// MonkeyAI makes random legal moves (it's a smart monkey!).
type MonkeyAI struct {
	b *board.KulamiBoard
}

// NewMonkeyAI creates a new monkey AI.
func NewMonkeyAI(b *board.KulamiBoard) *MonkeyAI {
	return &MonkeyAI{b}
}

// SuggestMove returns the best move this AI can come up with.
func (a *MonkeyAI) SuggestMove() (board.Coord, error) {
	moves := a.b.LegalMoves()
	if len(moves) == 0 {
		return board.Coord{}, ErrNoLegalMoves
	}
	return moves[rand.Intn(len(moves))], nil
}
