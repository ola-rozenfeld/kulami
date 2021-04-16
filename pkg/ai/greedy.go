package ai

import (
	"math"
	"math/rand"

	"github.com/ola-rozenfeld/kulami/pkg/board"
)

// GreedyAI takes the move maximizing the immediate score, with no look-ahead.
type GreedyAI struct {
	b *board.KulamiBoard
}

// NewGreedyAI creates a new greedy AI.
func NewGreedyAI(b *board.KulamiBoard) *GreedyAI {
	return &GreedyAI{b}
}

// SuggestMove returns the best move this AI can come up with.
func (a *GreedyAI) SuggestMove() (board.Coord, error) {
	moves := a.b.LegalMoves()
	if len(moves) == 0 {
		return board.Coord{}, ErrNoLegalMoves
	}
	b := a.b.Clone()
	isRed := b.IsRedsTurn()
	var bestMoves []board.Coord
	bestValue := math.MinInt32
	for _, m := range moves {
		if err := b.Move(m, isRed); err != nil {
			return board.Coord{}, err
		}
		score := b.ScoreDiff(isRed)
		b.UndoLastMove()
		if score == bestValue {
			bestMoves = append(bestMoves, m)
		} else if score > bestValue {
			bestValue = score
			bestMoves = []board.Coord{m}
		}
	}
	return bestMoves[rand.Intn(len(bestMoves))], nil
}
