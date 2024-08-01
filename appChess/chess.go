package appChess

import (
	"strings"

	"github.com/neutralusername/Systemge/Tools"
)

func (game *ChessGame) marshalBoard() string {
	var builder strings.Builder
	for _, row := range game.board {
		for _, piece := range row {
			if piece == nil {
				builder.WriteString(".")
			} else {
				builder.WriteString(piece.getLetter())
			}
		}
	}
	return builder.String()
}

func (game *ChessGame) isWhiteTurn() bool {
	return len(game.moves)%2 == 0
}

func getStandardStartingPosition() [8][8]Piece {
	return [8][8]Piece{
		{&Rook{true, false}, &Knight{true}, &Bishop{true}, &Queen{true}, &King{true, false}, &Bishop{true}, &Knight{true}, &Rook{true, false}},
		{&Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{&Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}},
		{&Rook{false, false}, &Knight{false}, &Bishop{false}, &Queen{false}, &King{false, false}, &Bishop{false}, &Knight{false}, &Rook{false, false}},
	}
}

func get960StartingPosition() [8][8]Piece {
	randomizer := Tools.NewRandomizer(Tools.GetSystemTime())

	bishop1 := randomizer.GenerateRandomNumber(0, 3) * 2
	bishop2 := randomizer.GenerateRandomNumber(0, 3)*2 + 1

	pieces := []int{0, 1, 2, 3, 4, 5, 6, 7}
	remove(&pieces, int(bishop1))
	remove(&pieces, int(bishop2))
	for i := range pieces {
		j := randomizer.GenerateRandomNumber(0, int64(i))
		pieces[i], pieces[j] = pieces[j], pieces[i]
	}

	remaining := pieces[:5]
	rook1 := remaining[0]
	rook2 := remaining[1]
	if rook1 > rook2 {
		rook1, rook2 = rook2, rook1
	}
	king := remaining[2]
	if king < rook1 || king > rook2 {
		king = remaining[3]
		if king < rook1 || king > rook2 {
			king = remaining[4]
		}
	}

	finalPositions := make([]int, 8)
	finalPositions[bishop1] = 1
	finalPositions[bishop2] = 1
	finalPositions[rook1] = 0
	finalPositions[rook2] = 0
	finalPositions[king] = 4

	knightsPlaced := 0
	queenPlaced := false
	for i := 0; i < 8; i++ {
		if finalPositions[i] == 0 {
			if knightsPlaced < 2 {
				finalPositions[i] = 2
				knightsPlaced++
			} else if !queenPlaced {
				finalPositions[i] = 3
				queenPlaced = true
			} else {
				finalPositions[i] = 0
			}
		}
	}

	var position [8][8]Piece
	for i, piece := range finalPositions {
		switch piece {
		case 0:
			position[0][i] = &Rook{true, false}
			position[7][i] = &Rook{false, false}
		case 1:
			position[0][i] = &Bishop{true}
			position[7][i] = &Bishop{false}
		case 2:
			position[0][i] = &Knight{true}
			position[7][i] = &Knight{false}
		case 3:
			position[0][i] = &Queen{true}
			position[7][i] = &Queen{false}
		case 4:
			position[0][i] = &King{true, false}
			position[7][i] = &King{false, false}
		}
	}

	for i := 0; i < 8; i++ {
		position[1][i] = &Pawn{true}
		position[6][i] = &Pawn{false}
	}

	return position
}

func remove(slice *[]int, s int) {
	for i, v := range *slice {
		if v == s {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}
