
export class chessBoard extends React.Component {
	constructor(props) {
		super(props)
	}

    render() {
        let board = []
        for (let i = 0; i < 8; i++) {
            let row = []
            for (let j = 0; j < 8; j++) {
                let square = this.props.board[i*8 + j]
                row.push(React.createElement("div", {
                        style : {
                            width : "50px",
                            height : "50px",
                            backgroundColor : (i + j) % 2 === 0 ? "white" : "black",
                            display : "flex",
                            alignItems : "center",
                            justifyContent : "center",
                            color : "red",
                            outline : this.props.selected && this.props.selected.row === i && this.props.selected.col === j ? "5px solid #7CB9E8" : "none",
                        },
                        onClick : () => {
                            if (this.props.selected === null) {
                                    this.props.setStateRoot({
                                        selected : {
                                            row : i,
                                            col : j,
                                        }
                                    })
                            } else {
                                if (this.props.selected.row === i && this.props.selected.col === j) {
                                    this.props.setStateRoot({
                                        selected : null,
                                    
                                    })
                                    return
                                }
                                this.props.WS_CONNECTION.send(this.props.constructMessage("move", JSON.stringify({
                                    fromRow : this.props.selected.row,
                                    fromCol : this.props.selected.col,
                                    toRow : i,
                                    toCol : j,
                                })))
                                this.props.setStateRoot({
                                    selected : null,
                                })
                            }
                        },
                    },
                    square
                ))
            }
            board.push(React.createElement("div", {
                style : {
                    display : "flex",
                    flexDirection : "row",
                }
            }, row))
        }
        return React.createElement("div", {
				style : {   
                    display : "flex",
                    flexDirection : "column",
                    alignItems : "center",
                    justifyContent : "center",
				}
			},
            board
		)
    }
}