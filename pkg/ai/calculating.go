package ai

import (
	"errors"
	"fmt"

	"github.com/ola-rozenfeld/kulami/pkg/board"
)

var ErrNoLegalMoves = errors.New("no legal moves")

// KulamiAI knows how to play the game.
type KulamiAI interface {
	SuggestMove() (board.Coord, error)
}

// CalculatingAI searches for the best move.
type CalculatingAI struct {
	b *board.KulamiBoard
}

// NewCalculatingAI creates a new greedy AI.
func NewCalculatingAI(b *board.KulamiBoard) *CalculatingAI {
	return &CalculatingAI{b}
}

// SuggestMove returns the best move this AI can come up with.
func (a *CalculatingAI) SuggestMove() (board.Coord, error) {
	return board.Coord{}, fmt.Errorf("not implemented yet")
}
