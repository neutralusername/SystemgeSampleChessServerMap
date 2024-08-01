import { game } from "./game.js";
import { home } from "./home.js";

export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
                id : "",
                idInput : "",
                errorMessage : "",
                errorTimeout : null,
                content: "",
                board: "",
                moves : [],
                selected: null,
                WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),
                constructMessage: (topic, payload) => {
                    return JSON.stringify({
                        topic: topic,
                        payload: payload,
                    });
                },
                setStateRoot: (state) => {
                    this.setState(state)
                },
                setErrorMessage: (message) => {
                    clearTimeout(this.state.errorTimeout)
                    this.setState({
                        errorMessage : message,
                        errorTimeout : setTimeout(() => {
                            this.setState({
                                errorMessage : "",
                            })
                        }, 5000)
                    })
                }
            },
            (this.state.WS_CONNECTION.onmessage = (event) => {
                let message = JSON.parse(event.data);
                console.log(message)
                switch (message.topic) {
                    case "connected":
                        this.state.setStateRoot({
                            id: message.payload,
                            content: "home",
                        });
                        break;
                    case "startGame":
                        this.state.setStateRoot({
                            content: "game",
                            board: message.payload,
                        });
                        break;
                    case "endGame":
                        this.state.setStateRoot({
                            content: "home",
                        });
                        break;
                    case "move":
                        let chessMove = JSON.parse(message.payload)
                        let fromIndex = chessMove.fromRow * 8 + chessMove.fromCol;
                        let boardCharArray = this.state.board.split("");
                        boardCharArray[chessMove.toRow * 8 + chessMove.toCol] = boardCharArray[fromIndex];
                        boardCharArray[fromIndex] = ".";
                        //abusing the setErrorMessage function to display the move
                        this.state.setErrorMessage((this.state.moves.length % 2 === 0 ? "white move: " : "black move: ") + chessMove.algebraicNotation)
                        this.state.setStateRoot({
                            board: boardCharArray.join(""),
                            moves: [...this.state.moves, chessMove],
                        });
                        break;
                    case "error":
                        let errorMessage = message.payload.split("->").reverse()[0]
                        this.state.setErrorMessage(errorMessage);
                        break;
                    default:
                        console.log("Unknown message topic: " + event.data);
                        break;
                }
            });
        this.state.WS_CONNECTION.onclose = () => {
            setTimeout(() => {
                if (this.state.WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                window.location.reload();
            }, 2000);
        };
        this.state.WS_CONNECTION.onopen = () => {
            let myLoop = () => {
                this.state.WS_CONNECTION.send(this.state.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

    render() {
        let content = null; 
        switch (this.state.content) {
            case "home":
                content = React.createElement(home, this.state);
                break;
            case "game":
                content = React.createElement(game, this.state);
                break;
            default:
                break;
        }
        return React.createElement(
            "div", {
                id: "root",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                },
            },
            React.createElement("div", {
                style: {
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                },
            }, this.state.errorMessage || "\u00a0"),
            content
        );
    }
}