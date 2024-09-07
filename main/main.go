package main

import (
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
)

const LOGGER_PATH = "logs.log"

func main() {
	Dashboard.NewServer("dashboardServer",
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
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
						Port:        60000,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{},
			},
			HeapUpdateIntervalMs:      1000,
			GoroutineUpdateIntervalMs: 1000,
			StatusUpdateIntervalMs:    1000,
			MetricsUpdateIntervalMs:   1000,
		},
	).Start()
	appWebsocketHTTP.New()
	appChess.New()
	<-make(chan time.Time)
}
