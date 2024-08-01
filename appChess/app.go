package appChess

import (
	"SystemgeSampleChessServer/dto"
	"sync"

	"github.com/neutralusername/Systemge/Node"
)

func newChessGame(whiteId string, blackId string) *ChessGame {
	game := &ChessGame{
		whiteId: whiteId,
		blackId: blackId,
	}
	game.initBoard()
	return game
}

func (chessGame *ChessGame) initBoard() {
	chessGame.board = getStandardStartingPosition()
}

type ChessGame struct {
	board   [8][8]Piece
	blackId string
	whiteId string
	moves   []*dto.Move
}

type App struct {
	games map[string]*ChessGame
	mutex sync.Mutex
}

func New() Node.Application {
	app := &App{
		games: make(map[string]*ChessGame),
	}
	return app
}
