package main

import (
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardServer"
)

const LOGGER_PATH = "logs.log"

func main() {
	DashboardServer.New("dashboardServer",
		&Config.DashboardServer{
			HTTPServerConfig: &Config.HTTPServer{
				TcpServerConfig: &Config.TcpServer{
					Port: 8081,
				},
			},
			WebsocketServerConfig: &Config.WebsocketServer{
				Pattern:                 "/ws",
				ClientWatchdogTimeoutMs: 1000 * 60,
				TcpServerConfig: &Config.TcpServer{
					Port: 8444,
				},
			},
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
						Port:        60000,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			},
			DashboardSystemgeCommands:   true,
			DashboardHttpCommands:       true,
			DashboardWebsocketCommands:  true,
			FrontendHeartbeatIntervalMs: 1000 * 60 * 1,
			UpdateIntervalMs:            1000,
			MaxEntriesPerMetrics:        100,
		},
		nil, nil,
	).Start()
	appWebsocketHTTP.New()
	appChess.New()
	<-make(chan time.Time)
}
