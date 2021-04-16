package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ola-rozenfeld/kulami/pkg/ai"
	"github.com/ola-rozenfeld/kulami/pkg/board"
)

// AIType is the supported levels of AI.
type AIType string

const (
	monkey      AIType = "monkey"
	greedy      AIType = "greedy"
	calculating AIType = "calculating"
)

var aiTypes = []AIType{monkey, greedy, calculating}

var (
	aiOpp  = flag.Bool("ai_opp", true, "Whether to play vs. an AI opponent or hot-seat.")
	aiType = flag.String("ai_type", string(monkey), fmt.Sprintf("Type/level of opponent AI. Supported values: %v", aiTypes))
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
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
	round := 1
	var aiEngine ai.KulamiAI
	aiPlayer := 1 //rand.Intn(2)
	if *aiOpp {
		fmt.Printf("Playing vs. the %s AI. The AI opponent is playing %s.\n", *aiType, playerNames[aiPlayer])
		switch *aiType {
		case string(monkey):
			aiEngine = ai.NewMonkeyAI(b)
		case string(greedy):
			aiEngine = ai.NewGreedyAI(b)
		case string(calculating):
			aiEngine = ai.NewCalculatingAI(b)
		}
	}
	for {
		fmt.Printf("%s", b)
		fmt.Printf("Round %d: it is %s to move. ", round, playerNames[player])
		var move board.Coord
		var err error
		if *aiOpp && player == aiPlayer {
			if move, err = aiEngine.SuggestMove(); err != nil {
				log.Fatalf("An AI error: %v", err)
			}
			fmt.Printf("AI chooses %d,%d.\n", move.Row, move.Col)
		} else {
			fmt.Printf("Type `resign` to resign, or a move coordinate: ")
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
			move = board.Coord{Row: row, Col: col}
		}
		if err := b.Move(move, player == 0); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		player = 1 - player
		if player == 0 {
			round++
		}
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
