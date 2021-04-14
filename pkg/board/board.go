// A Kulami board game engine.
package board

import (
	"fmt"
	"strings"
)

// Kulami has 17 tiles of the following sizes:
var kTileSizes = []int{6, 6, 6, 6, 4, 4, 4, 4, 4, 3, 3, 3, 3, 2, 2, 2, 2}

const (
	kNumTiles   = 17
	kNumMarbles = 28
	// Possible values of a coordinate on a board.
	kOutOfBounds = -1
	kEmptySpace  = iota
	kRedMarble
	kBlackMarble
)

// Coord is a point in 2D.
type Coord struct {
	Col, Row int
}

// KulamiBoard represents a full state in a Kulami game.
type KulamiBoard struct {
	end        Coord   // The lower-right corner of the board.
	tiles      [][]int // Indices of tile by coordinate.
	marbles    [][]int // Which marble, if any, exists at coordinate.
	moves      []Coord // All moves made thus far.
	redScore   int     // Total tiles with red majority so far.
	blackScore int     // Total tiles with blackMajority so far.
	tileScore  []int   // Marble advantage for red per tile.
}

// RedScore returns the current score of the red player.
func (b *KulamiBoard) RedScore() int {
	return b.redScore
}

// RedScore returns the current score of the black player.
func (b *KulamiBoard) BlackScore() int {
	return b.blackScore
}

// TileLocation represents a complete location of a tile inside a board.
type TileLocation struct {
	Coord       Coord
	IsLandscape bool
}

func (l TileLocation) tileEnd(t int) Coord {
	end := l.Coord
	switch kTileSizes[t] {
	case 6:
		if l.IsLandscape {
			end.Row += 1
			end.Col += 2
		} else {
			end.Row += 2
			end.Col += 1
		}
	case 4:
		end.Row += 1
		end.Col += 1
	case 3:
		if l.IsLandscape {
			end.Col += 2
		} else {
			end.Row += 2
		}
	case 2:
		if l.IsLandscape {
			end.Col += 1
		} else {
			end.Row += 1
		}
	}
	return end
}

// New initializes an empty Kulami board from tile coordinates.
// There should be exactly kNumPieces coordinates corresponding to the
// upper left corner of each tile.
func New(locs []TileLocation) (*KulamiBoard, error) {
	if len(locs) != kNumTiles {
		return nil, fmt.Errorf("need exactly %v tile locations, got %v", kNumTiles, len(locs))
	}
	b := &KulamiBoard{}
	// Compute the board range.
	for t, l := range locs {
		end := l.tileEnd(t)
		if end.Row > b.end.Row {
			b.end.Row = end.Row
		}
		if end.Col > b.end.Col {
			b.end.Col = end.Col
		}
	}
	b.tileScore = make([]int, kNumTiles)
	b.marbles = make([][]int, b.end.Row+1)
	b.tiles = make([][]int, b.end.Row+1)
	for i := range b.marbles {
		b.marbles[i] = make([]int, b.end.Col+1)
		b.tiles[i] = make([]int, b.end.Col+1)
		for j := range b.marbles[i] {
			b.marbles[i][j] = kOutOfBounds
			b.tiles[i][j] = kOutOfBounds
		}
	}
	for t, l := range locs {
		end := l.tileEnd(t)
		for row := l.Coord.Row; row <= end.Row; row++ {
			for col := l.Coord.Col; col <= end.Col; col++ {
				b.marbles[row][col] = kEmptySpace
				if b.tiles[row][col] != kOutOfBounds {
					return nil, fmt.Errorf("tiles %d and %d intersect on %d,%d", t, b.tiles[row][col], row, col)
				}
				b.tiles[row][col] = t
			}
		}
	}
	return b, nil
}

// String represenation of the board in the following format:
//
// *     0   1   2   3
//         ---------
// 0       | o   . |
//         |       -----
// 1       | .   . | o |
//     ------------|   |
// 2   | .   x   x | o |
//     ------------|   |
// 3               | . |
//                 -----
func (b *KulamiBoard) String() string {
	m1, m2 := Coord{Row: -1, Col: -1}, Coord{Row: -1, Col: -1}
	if len(b.moves) > 0 {
		m1 = b.moves[len(b.moves)-1]
	}
	if len(b.moves) > 1 {
		m2 = b.moves[len(b.moves)-2]
	}
	var res strings.Builder
	// Print column index.
	fmt.Fprint(&res, "*  ")
	for i := range b.marbles[0] {
		fmt.Fprintf(&res, "%4d", i)
	}
	fmt.Fprint(&res, "\n")
	for row, cols := range b.marbles {
		// A row of separators between rows.
		fmt.Fprint(&res, "    ")
		for col := range cols {
			if row == 0 && b.tiles[0][col] != kOutOfBounds || row != 0 && b.tiles[row][col] != b.tiles[row-1][col] {
				fmt.Fprint(&res, "----")
			} else if col == 0 && b.tiles[row][0] != kOutOfBounds || col != 0 && b.tiles[row][col] != b.tiles[row][col-1] {
				if b.tiles[row][col] == kOutOfBounds && (row == 0 || b.tiles[row-1][col-1] == kOutOfBounds) {
					fmt.Fprint(&res, "-   ")
				} else {
					fmt.Fprint(&res, "|   ")
				}
			} else if col != 0 && row != 0 && b.tiles[row][col] == kOutOfBounds && b.tiles[row-1][col-1] != kOutOfBounds {
				fmt.Fprint(&res, "-   ")
			} else {
				fmt.Fprint(&res, "    ")
			}
		}
		// Close the row from the right with either | or -.
		if b.tiles[row][b.end.Col] != kOutOfBounds || row != 0 && b.tiles[row-1][b.end.Col] != kOutOfBounds {
			// | is only for if we're inside the same tile.
			if row != 0 && b.tiles[row][b.end.Col] == b.tiles[row-1][b.end.Col] {
				fmt.Fprint(&res, "|")
			} else {
				fmt.Fprint(&res, "-")
			}
		}
		fmt.Fprint(&res, "\n")
		// Print row index.
		fmt.Fprintf(&res, "%-4d", row)
		// Print actual content (including marbles).
		for col, v := range cols {
			isLast := (row == m1.Row && col == m1.Col || row == m2.Row && col == m2.Col)
			m := " "
			switch v {
			case kEmptySpace:
				m = "."
			case kRedMarble:
				m = "x"
				if isLast {
					m = "X"
				}
			case kBlackMarble:
				m = "o"
				if isLast {
					m = "O"
				}
			}
			sep := " "
			if col == 0 && b.tiles[row][0] != kOutOfBounds || col != 0 && b.tiles[row][col] != b.tiles[row][col-1] {
				sep = "|"
			}
			fmt.Fprintf(&res, "%s %s ", sep, m)
		}
		// Possibly close the row from the right with |.
		if b.tiles[row][b.end.Col] != kOutOfBounds {
			fmt.Fprint(&res, "|")
		}
		fmt.Fprint(&res, "\n")
	}
	// Last closing row.
	fmt.Fprint(&res, "    ")
	for col, t := range b.tiles[b.end.Row] {
		if t != kOutOfBounds {
			fmt.Fprint(&res, "----")
		} else if col != 0 && b.tiles[b.end.Row][col-1] != kOutOfBounds {
			fmt.Fprint(&res, "-   ")
		} else {
			fmt.Fprint(&res, "    ")
		}
	}
	fmt.Fprint(&res, "\n")
	fmt.Fprintf(&res, "\nScore:\t\tRed: %d\tBlack: %d\n", b.redScore, b.blackScore)
	return res.String()
}

// Apply a move to the board, if legal.
func (b *KulamiBoard) Move(c Coord, isRed bool) error {
	last := Coord{Row: -1, Col: -1}
	if len(b.moves) != 0 {
		last = b.moves[len(b.moves)-1]
	}
	if last.Row >= 0 && isRed == (b.marbles[last.Row][last.Col] == kRedMarble) {
		return fmt.Errorf("it is now the other player's turn")
	}
	if len(b.moves) == kNumMarbles*2 {
		return fmt.Errorf("game is over, out of marbles")
	}
	if c.Row < 0 || c.Row > b.end.Row || c.Col < 0 || c.Col > b.end.Col {
		return fmt.Errorf("%d,%d is not a legal move", c.Row, c.Col)
	}
	if b.marbles[c.Row][c.Col] != kEmptySpace || last.Row >= 0 && last.Row != c.Row && last.Col != c.Col {
		return fmt.Errorf("%d,%d is not a legal move", c.Row, c.Col)
	}
	if last.Row >= 0 && b.tiles[c.Row][c.Col] == b.tiles[last.Row][last.Col] {
		return fmt.Errorf("%d,%d is not a legal move, the tile is blocked by the other player", c.Row, c.Col)
	}
	if len(b.moves) > 1 {
		pre := b.moves[len(b.moves)-2]
		if b.tiles[c.Row][c.Col] == b.tiles[pre.Row][pre.Col] {
			return fmt.Errorf("%d,%d is not a legal move, the tile is blocked by you", c.Row, c.Col)
		}
	}
	b.moves = append(b.moves, c)
	if isRed {
		b.marbles[c.Row][c.Col] = kRedMarble
	} else {
		b.marbles[c.Row][c.Col] = kBlackMarble
	}
	tile := b.tiles[c.Row][c.Col]
	curScore := b.tileScore[tile]
	delta := 0
	if curScore == 0 || curScore == -1 && isRed || curScore == 1 && !isRed {
		delta = kTileSizes[tile]
		if isRed {
			if curScore == 0 {
				b.redScore += delta
			} else {
				b.blackScore -= delta
			}
		} else {
			if curScore == 0 {
				b.blackScore += delta
			} else {
				b.redScore -= delta
			}
		}
	}
	if isRed {
		b.tileScore[tile] += 1
	} else {
		b.tileScore[tile] -= 1
	}
	return nil
}

// UndoLastMove removes the last move from the board, if possible.
func (b *KulamiBoard) UndoLastMove() {
	if len(b.moves) == 0 {
		return
	}
	c := b.moves[len(b.moves)-1]
	b.moves = b.moves[0 : len(b.moves)-1]
	isRed := b.marbles[c.Row][c.Col] == kRedMarble
	b.marbles[c.Row][c.Col] = kEmptySpace
	tile := b.tiles[c.Row][c.Col]
	curScore := b.tileScore[tile]
	delta := 0
	if curScore == 0 || curScore == 1 && isRed || curScore == -1 && !isRed {
		delta = kTileSizes[tile]
		if isRed {
			if curScore == 0 {
				b.blackScore += delta
			} else {
				b.redScore -= delta
			}
		} else {
			if curScore == 0 {
				b.redScore += delta
			} else {
				b.blackScore -= delta
			}
		}
	}
	if isRed {
		b.tileScore[tile] -= 1
	} else {
		b.tileScore[tile] += 1
	}
}

// LegalMoves returns all legal move candidates for the next move.
func (b *KulamiBoard) LegalMoves() []Coord {
	var res []Coord
	if len(b.moves) == 0 {
		// All moves on the board are legal as a first move.
		for row, cols := range b.marbles {
			for col, v := range cols {
				if v == kEmptySpace {
					res = append(res, Coord{Row: row, Col: col})
				}
			}
		}
		return res
	}
	last := b.moves[len(b.moves)-1]
	if len(b.moves) == kNumMarbles*2 {
		return res // Game over, out of marbles.
	}
	prev := Coord{Row: -1, Col: -1}
	if len(b.moves) > 1 {
		prev = b.moves[len(b.moves)-2]
	}
	for row := 0; row <= b.end.Row; row++ {
		if b.marbles[row][last.Col] == kEmptySpace &&
			b.tiles[row][last.Col] != b.tiles[last.Row][last.Col] &&
			(prev.Row == -1 || b.tiles[row][last.Col] != b.tiles[prev.Row][prev.Col]) {
			res = append(res, Coord{Row: row, Col: last.Col})
		}
	}
	for col := 0; col <= b.end.Col; col++ {
		if b.marbles[last.Row][col] == kEmptySpace &&
			b.tiles[last.Row][col] != b.tiles[last.Row][last.Col] &&
			(prev.Row == -1 || b.tiles[last.Row][col] != b.tiles[prev.Row][prev.Col]) {
			res = append(res, Coord{Row: last.Row, Col: col})
		}
	}
	return res
}
