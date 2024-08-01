package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/dto"
	"strings"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"startGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !node.WebsocketClientExists(blackId) {
				return Error.New("Opponent does not exist", nil)
			}
			responseChannel, err := node.SyncMessage("startGame", Helpers.JsonMarshal([]string{whiteId, blackId}))
			if err != nil {
				return Error.New("Error sending start message", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving start response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			err = node.AddToWebsocketGroup(whiteId+"-"+blackId, whiteId, blackId)
			if err != nil {
				responseChannel, err := node.SyncMessage("endGame", whiteId+"-"+blackId)
				if err != nil {
					if errorLogger := node.GetErrorLogger(); errorLogger != nil {
						errorLogger.Log(Error.New("Error sending endGame message", err).Error())
					}
				}
				response, err := responseChannel.ReceiveResponse()
				if err != nil {
					if errorLogger := node.GetErrorLogger(); errorLogger != nil {
						errorLogger.Log(Error.New("Error receiving endGame response", err).Error())
					}
				}
				if response.GetTopic() == Message.TOPIC_FAILURE {
					if errorLogger := node.GetErrorLogger(); errorLogger != nil {
						errorLogger.Log(Error.New(response.GetPayload(), nil).Error())
					}
				}
				return Error.New("Error adding to group", err)
			}
			node.WebsocketGroupcast(whiteId+"-"+blackId, Message.NewAsync("propagate_gameStart", response.GetPayload()))
			return nil
		},
		"endGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			responseChannel, err := node.SyncMessage("endGame", websocketClient.GetId())
			if err != nil {
				return Error.New("Error sending endGame message", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving endGame response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			gameId := response.GetPayload()
			ids := strings.Split(gameId, "-")
			node.WebsocketGroupcast(gameId, Message.NewAsync("propagate_gameEnd", ""))
			node.RemoveFromWebsocketGroup(gameId, ids...)
			return nil
		},
		"move": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			move, err := dto.UnmarshalMove(message.GetPayload())
			if err != nil {
				return Error.New("Error unmarshalling move", err)
			}
			move.PlayerId = websocketClient.GetId()
			responseChannel, err := node.SyncMessage("move", Helpers.JsonMarshal(move))
			if err != nil {
				return Error.New("Error sending move message", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving move response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			responseMove, err := dto.UnmarshalMove(response.GetPayload())
			if err != nil {
				return Error.New("Error unmarshalling response move", err)
			}
			node.WebsocketGroupcast(responseMove.GameId, Message.NewAsync("propagate_move", Helpers.JsonMarshal(responseMove)))
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", websocketClient.GetId()).Serialize())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending connected message", err).Error())
		}
		websocketClient.Disconnect()
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	responseChannel, err := node.SyncMessage("endGame", websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending endGame message", err).Error())
		}
	}
	response, err := responseChannel.ReceiveResponse()
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error receiving endGame response", err).Error())
		}
	}
	if response.GetTopic() == Message.TOPIC_SUCCESS {
		gameId := response.GetPayload()
		ids := strings.Split(gameId, "-")
		node.WebsocketGroupcast(gameId, Message.NewAsync("propagate_gameEnd", ""))
		node.RemoveFromWebsocketGroup(gameId, ids...)
	}
}
