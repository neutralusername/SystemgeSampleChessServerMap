package appChess

import (
	"sync"

	"github.com/neutralusername/Systemge/Node"
)

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
