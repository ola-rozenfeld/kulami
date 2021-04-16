package board

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var sampleTiles = []TileLocation{
	// 6s
	{Coord: Coord{Row: 4, Col: 0}},
	{Coord: Coord{Row: 6, Col: 2}, IsLandscape: true},
	{Coord: Coord{Row: 4, Col: 3}, IsLandscape: true},
	{Coord: Coord{Row: 1, Col: 6}},
	// 4s
	{Coord: Coord{Row: 0, Col: 4}},
	{Coord: Coord{Row: 2, Col: 4}},
	{Coord: Coord{Row: 2, Col: 2}},
	{Coord: Coord{Row: 4, Col: 6}},
	{Coord: Coord{Row: 7, Col: 5}},
	// 3s
	{Coord: Coord{Row: 1, Col: 1}, IsLandscape: true},
	{Coord: Coord{Row: 2, Col: 8}},
	{Coord: Coord{Row: 5, Col: 8}, IsLandscape: true},
	{Coord: Coord{Row: 6, Col: 5}, IsLandscape: true},
	//2s
	{Coord: Coord{Row: 4, Col: 2}},
	{Coord: Coord{Row: 4, Col: 9}, IsLandscape: true},
	{Coord: Coord{Row: 2, Col: 0}},
	{Coord: Coord{Row: 2, Col: 1}},
}
var sampleMoves = []Coord{
	{Row: 4, Col: 5},
	{Row: 4, Col: 0},
	{Row: 2, Col: 0},
	{Row: 2, Col: 7},
	{Row: 4, Col: 7},
	{Row: 4, Col: 3},
	{Row: 4, Col: 1},
	{Row: 4, Col: 6},
	{Row: 7, Col: 6},
	{Row: 2, Col: 6},
	{Row: 6, Col: 6},
	{Row: 5, Col: 6},
	{Row: 1, Col: 6},
}

const printOut = `
*     0   1   2   3   4   5   6   7   8   9  10
                    ---------                   
0                   | .   . |                   
        ------------|       ---------           
1       | .   .   . | .   . | X   . |           
    ------------------------|       -----       
2   | x | . | .   . | .   . | o   o | . |       
    |   |   |       |       |       |   |       
3   | . | . | .   . | .   . | .   . | . |       
    --------------------------------|   ---------
4   | o   x | . | o   .   x | o   x | . | .   . |
    |       |   |           |       -------------
5   | .   . | . | .   .   . | O   . | .   .   . |
    |       -------------------------------------
6   | .   . | .   .   . | .   x   . |           
    --------|           -------------           
7           | .   .   . | .   x |               
            ------------|       |               
8                       | .   . |               
                        ---------               

Score:		Red: 9	Black: 10
`

func TestMove(t *testing.T) {
	b, err := New(sampleTiles)
	if err != nil {
		t.Fatalf("Error initializing board: %v", err)
	}
	isRed := true
	for _, m := range sampleMoves {
		if err := b.Move(m, isRed); err != nil {
			t.Errorf("Move(%d,%d): %v", m.Row, m.Col, err)
		}
		isRed = !isRed
	}
	if got, want := b.RedScore(), 9; got != want {
		t.Errorf("RedScore() = %d, want %d", got, want)
	}
	if got, want := b.BlackScore(), 10; got != want {
		t.Errorf("BlackScore() = %d, want %d", got, want)
	}
	// Check Undo works.
	for range sampleMoves {
		b.UndoLastMove()
	}
	if got, want := b.RedScore(), 0; got != want {
		t.Errorf("RedScore() = %d, want %d", got, want)
	}
	if got, want := b.BlackScore(), 0; got != want {
		t.Errorf("BlackScore() = %d, want %d", got, want)
	}
}

func TestMoveErrors(t *testing.T) {
	b, err := New(sampleTiles)
	if err != nil {
		t.Fatalf("Error initializing board: %v", err)
	}
	b.Move(Coord{Row: 4, Col: 5}, true)
	b.Move(Coord{Row: 4, Col: 0}, false)
	illegalMoves := []Coord{
		{Row: 4, Col: 5},
		{Row: 4, Col: 0},
		{Row: 4, Col: 3},
		{Row: 5, Col: 0},
		{Row: 7, Col: 0},
		{Row: 4, Col: 11},
	}
	for _, m := range illegalMoves {
		if b.Move(m, true) == nil {
			t.Errorf("Expected Move(%d,%d) to error", m.Row, m.Col)
		}
	}
	if b.Move(Coord{Row: 2, Col: 0}, false) == nil {
		t.Errorf("Expected Move(2,0) to error") // Red's turn.
	}
}

func TestLegalMoves(t *testing.T) {
	b, err := New(sampleTiles)
	if err != nil {
		t.Fatalf("Error initializing board: %v", err)
	}
	isRed := true
	for _, m := range sampleMoves {
		if err := b.Move(m, isRed); err != nil {
			t.Errorf("Move(%d,%d): %v", m.Row, m.Col, err)
		}
		isRed = !isRed
	}
	wantMoves := []Coord{
		{Row: 8, Col: 6},
		{Row: 1, Col: 1},
		{Row: 1, Col: 2},
		{Row: 1, Col: 3},
		{Row: 1, Col: 4},
		{Row: 1, Col: 5},
	}
	gotMoves := b.LegalMoves()
	if diff := cmp.Diff(wantMoves, gotMoves); diff != "" {
		t.Errorf("LegalMoves mismatch (-want +got):\n%s", diff)
	}
}

func TestString(t *testing.T) {
	b, err := New(sampleTiles)
	if err != nil {
		t.Fatalf("Error initializing board: %v", err)
	}
	isRed := true
	for _, m := range sampleMoves {
		if err := b.Move(m, isRed); err != nil {
			t.Errorf("Move(%d,%d): %v", m.Row, m.Col, err)
		}
		isRed = !isRed
	}
	if got := b.String(); "\n"+got != printOut {
		t.Errorf("String() returned:\n%s\nExpected:\n%s\n", got, printOut)
	}
}
