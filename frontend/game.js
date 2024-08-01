import { chessBoard } from "./chessBoard.js"

export class game extends React.Component {
	constructor(props) {
		super(props)
	}

    render() {
		let moves = []
		this.props.moves.forEach((move, i) => {
			let moveNumber = (i % 2 == 0) ? (i / 2 + 1) + ". " : (i + 1) / 2 + ". "
			moves.push(React.createElement("div", {
					style : {
						marginTop : "1vmin",
						fontSize : "2vmin",
					}
				},moveNumber+ move.algebraicNotation
			))
		})
        return React.createElement("div", {
				style : {
					gap : "1vmin",
					position: "relative",
					marginTop : "1vmin",
					display: "flex",
					flexDirection : "column",
					alignItems : "center",
					justifyContent : "center",
				}
			},
			React.createElement(chessBoard, this.props),
			React.createElement("button", {
					style : {
						marginTop : "1vmin",
					},
					onClick : () => {
						this.props.WS_CONNECTION.send(this.props.constructMessage("endGame", ""))
					}
				}, "End Game",
			),
			moves,
		)
    }
}