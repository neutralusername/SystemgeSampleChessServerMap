package appChess

import (
	"SystemgeSampleChessServer/dto"
	"SystemgeSampleChessServer/topics"
	"encoding/json"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SingleRequestServer"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type App struct {
	games map[string]*ChessGame
	mutex sync.Mutex

	singleRequestServer *SingleRequestServer.Server
}

func New() *App {
	app := &App{
		games: make(map[string]*ChessGame),
	}

	app.singleRequestServer = SingleRequestServer.NewSingleRequestServer("chessSpawner",
		&Config.SingleRequestServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
						Port:        60001,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{},
			},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpSystemgeConnection{},
				ClientConfig: &Config.TcpClient{
					Address: "localhost:60000",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
		},
		nil, SystemgeConnection.NewConcurrentMessageHandler(
			SystemgeConnection.AsyncMessageHandlers{},
			SystemgeConnection.SyncMessageHandlers{
				topics.MOVE: func(systemgeConnection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					app.mutex.Lock()
					defer app.mutex.Unlock()
					move, err := dto.UnmarshalMove(message.GetPayload())
					if err != nil {
						return "", Error.New("Error unmarshalling move", err)
					}
					game := app.games[move.PlayerId]
					if game == nil {
						return "", Error.New("Game does not exist", nil)
					}
					move.GameId = game.whiteId + "-" + game.blackId
					chessMove, err := game.handleMoveRequest(move)
					if err != nil {
						return "", err
					}
					return Helpers.JsonMarshal(chessMove), nil
				},
				topics.STARTGAME: func(systemgeConnection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					ids := []string{}
					json.Unmarshal([]byte(message.GetPayload()), &ids)
					whiteId := ids[0]
					blackId := ids[1]
					game := newChessGame(whiteId, blackId)
					app.mutex.Lock()
					defer app.mutex.Unlock()
					if app.games[whiteId] != nil || app.games[blackId] != nil {
						app.mutex.Unlock()
						return "", Error.New("Already in a game", nil)
					}
					app.games[whiteId] = game
					app.games[blackId] = game
					return game.marshalBoard(), nil
				},
				topics.ENDGAME: func(systemgeConnection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					id := message.GetPayload()
					app.mutex.Lock()
					defer app.mutex.Unlock()
					if app.games[id] == nil {
						return "", Error.New("Game does not exist", nil)
					}
					game := app.games[id]
					delete(app.games, game.whiteId)
					delete(app.games, game.blackId)
					return game.whiteId + "-" + game.blackId, nil
				},
			},
			nil, nil,
		),
	)
	if err := app.singleRequestServer.Start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start singleRequestServer", err))
	}

	return app
}

func (game *ChessGame) handleMoveRequest(move *dto.Move) (*dto.Move, error) {
	if game.isWhiteTurn() && move.PlayerId != game.whiteId {
		return nil, Error.New("Not your turn", nil)
	}
	if !game.isWhiteTurn() && move.PlayerId != game.blackId {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := game.move(move)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}
