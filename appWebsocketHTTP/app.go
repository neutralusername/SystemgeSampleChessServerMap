package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/dto"
	"SystemgeSampleChessServer/topics"
	"errors"
	"strings"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SingleRequestServer"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{
		status: Status.STOPPED,
	}

	app.websocketServer = WebsocketServer.New("appWebsocketHttp_websocketServer",
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port: 8443,
			},
		},
		WebsocketServer.MessageHandlers{
			topics.STARTGAME: func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				whiteId := websocketClient.GetId()
				blackId := message.GetPayload()
				if !app.websocketServer.ClientExists(blackId) {
					return Error.New("Opponent does not exist", nil)
				}
				response, err := SingleRequestServer.SyncRequest("websocketHttpServer", &Config.SingleRequestClient{
					TcpConnectionConfig: &Config.TcpSystemgeConnection{},
					TcpClientConfig: &Config.TcpClient{
						Address: "localhost:60001",
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					},
				}, topics.STARTGAME, Helpers.JsonMarshal([]string{whiteId, blackId}))
				if err != nil {
					panic(Error.New("Error sending startGame message", err))
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					panic(Error.New("Error starting game", errors.New(response.GetPayload())))
				}
				err = app.websocketServer.AddClientToGroup(whiteId+"-"+blackId, whiteId, blackId)
				if err != nil {
					panic(Error.New("Error adding clients to group", err))
				}
				err = app.websocketServer.Groupcast(whiteId+"-"+blackId, Message.NewAsync(topics.STARTGAME, response.GetPayload()))
				if err != nil {
					panic(Error.New("Error groupcasting start message", err))
				}
				return nil
			},
			topics.ENDGAME: func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				response, err := SingleRequestServer.SyncRequest("websocketHttpServer", &Config.SingleRequestClient{
					TcpConnectionConfig: &Config.TcpSystemgeConnection{},
					TcpClientConfig: &Config.TcpClient{
						Address: "localhost:60001",
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					},
				}, topics.ENDGAME, websocketClient.GetId())
				if err != nil {
					return Error.New("Error sending endGame message", err)
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					return Error.New("Error ending game", nil)
				}
				ids := strings.Split(response.GetPayload(), "-")
				app.websocketServer.Groupcast(response.GetPayload(), Message.NewAsync(topics.ENDGAME, ""))
				app.websocketServer.RemoveClientFromGroup(response.GetPayload(), ids...)
				return nil
			},
			topics.MOVE: func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				move, err := dto.UnmarshalMove(message.GetPayload())
				if err != nil {
					return Error.New("Error unmarshalling move", err)
				}
				move.PlayerId = websocketClient.GetId()
				response, err := SingleRequestServer.SyncRequest("websocketHttpServer", &Config.SingleRequestClient{
					TcpConnectionConfig: &Config.TcpSystemgeConnection{},
					TcpClientConfig: &Config.TcpClient{
						Address: "localhost:60001",
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					},
				}, topics.MOVE, Helpers.JsonMarshal(move))
				if err != nil {
					return Error.New("Error sending move message", err)
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					return Error.New("Error making move", nil)
				}
				responseMove, err := dto.UnmarshalMove(response.GetPayload())
				if err != nil {
					return Error.New("Error unmarshalling response move", err)
				}
				app.websocketServer.Groupcast(responseMove.GameId, Message.NewAsync(topics.MOVE, Helpers.JsonMarshal(responseMove)))
				return nil
			},
		},
		app.OnConnectHandler, app.OnDisconnectHandler,
	)
	app.httpServer = HTTPServer.New("httpServer",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8080,
			},
		},
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)
	Dashboard.NewClient("appWebsocketHttp_dashboardClient",
		&Config.DashboardClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{},
			ClientConfig: &Config.TcpClient{
				Address: "localhost:60001",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.start, app.stop, nil, app.getStatus,
		nil,
	).Start()
	if err := app.start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start appWebsocketHttp", err))
	}
	return app
}

func (app *AppWebsocketHTTP) getStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.websocketServer.Start(); err != nil {
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.status = Status.STOPPED
	return nil
}

func (app *AppWebsocketHTTP) WebsocketPropagate(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	err := websocketClient.Send(Message.NewAsync("connected", websocketClient.GetId()).Serialize())
	if err != nil {
		return Error.New("Error sending connected message", err)
	}
	return nil
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketServer.WebsocketClient) {
	response, err := SingleRequestServer.SyncRequest("websocketHttpServer", &Config.SingleRequestClient{
		TcpConnectionConfig: &Config.TcpSystemgeConnection{},
		TcpClientConfig: &Config.TcpClient{
			Address: "localhost:60001",
			TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			Domain:  "example.com",
		},
	}, topics.ENDGAME, websocketClient.GetId())
	if err != nil {
		panic(Error.New("Error sending endGame request", err))
	}
	if response.GetTopic() == Message.TOPIC_SUCCESS {
		gameId := response.GetPayload()
		ids := strings.Split(gameId, "-")
		app.websocketServer.Groupcast(gameId, Message.NewAsync("propagate_gameEnd", ""))
		app.websocketServer.RemoveClientFromGroup(gameId, ids...)
	}
}
