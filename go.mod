module SystemgeSampleChessServer

go 1.23

toolchain go1.23.0

//replace github.com/neutralusername/Systemge => ../Systemge

require github.com/neutralusername/Systemge v0.0.0-20240907181301-726447715837

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
)
