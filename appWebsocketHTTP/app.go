package appWebsocketHTTP

import (
	"github.com/neutralusername/Systemge/Node"
)

type AppWebsocketHTTP struct {
}

func New() *AppWebsocketHTTP {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{}
}
