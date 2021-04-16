package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ola-rozenfeld/kulami/pkg/board"
)

func main() {
	sampleBoard := []board.TileLocation{
		// 6s
		{Coord: board.Coord{Row: 4, Col: 0}},
		{Coord: board.Coord{Row: 6, Col: 2}, IsLandscape: true},
		{Coord: board.Coord{Row: 4, Col: 3}, IsLandscape: true},
		{Coord: board.Coord{Row: 1, Col: 6}},
		// 4s
		{Coord: board.Coord{Row: 0, Col: 4}},
		{Coord: board.Coord{Row: 2, Col: 4}},
		{Coord: board.Coord{Row: 2, Col: 2}},
		{Coord: board.Coord{Row: 4, Col: 6}},
		{Coord: board.Coord{Row: 7, Col: 5}},
		// 3s
		{Coord: board.Coord{Row: 1, Col: 1}, IsLandscape: true},
		{Coord: board.Coord{Row: 2, Col: 8}},
		{Coord: board.Coord{Row: 5, Col: 8}, IsLandscape: true},
		{Coord: board.Coord{Row: 6, Col: 5}, IsLandscape: true},
		//2s
		{Coord: board.Coord{Row: 4, Col: 2}},
		{Coord: board.Coord{Row: 4, Col: 9}, IsLandscape: true},
		{Coord: board.Coord{Row: 2, Col: 0}},
		{Coord: board.Coord{Row: 2, Col: 1}},
	}
	b, err := board.New(sampleBoard)
	if err != nil {
		log.Fatalf("Error initializing board: %v", err)
	}
	player := 0
	playerNames := []string{"Red", "Black"}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s", b)
		fmt.Printf("It is %s to move. Type `resign` to resign, or a move coordinate: ", playerNames[player])
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "resign" {
			fmt.Printf("%s resigned. %s wins.\n", playerNames[player], playerNames[1-player])
			return
		}
		toks := strings.Split(strings.TrimSpace(text), ",")
		if len(toks) != 2 {
			fmt.Println("Error: expected coordinate as Row,Col")
			continue
		}
		row, err := strconv.Atoi(toks[0])
		if err != nil {
			fmt.Printf("Error %v: expected coordinate as Row,Col\n", err)
			continue
		}
		col, err := strconv.Atoi(toks[1])
		if err != nil {
			fmt.Printf("Error %v: expected coordinate as Row,Col\n", err)
			continue
		}
		if err := b.Move(board.Coord{Row: row, Col: col}, player == 0); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		player = 1 - player
		if len(b.LegalMoves()) == 0 {
			redScore := b.RedScore()
			blackScore := b.BlackScore()
			if redScore == blackScore {
				fmt.Printf("The game is a draw with final score %d vs. %d!\n", redScore, blackScore)
			} else if redScore > blackScore {
				fmt.Printf("Red wins with final score %d vs. %d!\n", redScore, blackScore)
			} else {
				fmt.Printf("Black wins with final score %d vs. %d!\n", blackScore, redScore)
			}
			return
		}
	}
}
