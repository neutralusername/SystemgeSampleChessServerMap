package appChess

import (
	"SystemgeSampleChessServer/dto"
	"SystemgeSampleChessServer/topics"
	"encoding/json"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.MOVE: func(node *Node.Node, message *Message.Message) (string, error) {
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
		topics.STARTGAME: func(node *Node.Node, message *Message.Message) (string, error) {
			ids := []string{}
			json.Unmarshal([]byte(message.GetPayload()), &ids)
			whiteId := ids[0]
			blackId := ids[1]
			game := newChessGame(whiteId, blackId)
			app.mutex.Lock()
			if app.games[whiteId] != nil || app.games[blackId] != nil {
				app.mutex.Unlock()
				return "", Error.New("Already in a game", nil)
			}
			app.games[whiteId] = game
			app.games[blackId] = game
			app.mutex.Unlock()
			return game.marshalBoard(), nil
		},
		topics.ENDGAME: func(node *Node.Node, message *Message.Message) (string, error) {
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
	}
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
