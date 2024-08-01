
export class home extends React.Component {
	constructor(props) {
		super(props)
	}

    render() {
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
            "your id: " + this.props.id,
            React.createElement("div", {
                    style: {
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                        alignItems: "center",
                    },
                },
                "enter another id or share your id"
            ),
            React.createElement("div", {
                    style: {
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "center",
                        alignItems: "center",
                    },
                },
                React.createElement("input", {
                    type: "text",
                    id: "input",
                    value: this.props.idInput,
                    onChange: (e) => {
                        this.props.setStateRoot({
                            idInput: e.target.value,
                        });
                    },
                    style: {
                        width: "100px",
                        height: "20px",
                    },
                }),
                React.createElement("button", {
                    onClick: () => {
                        this.props.WS_CONNECTION.send(
                           this.props.constructMessage("startGame", this.props.idInput)
                        );
                    },
                    style: {
                        width: "100px",
                        height: "20px",
                    },
                }, "start game")
            ),
		)
    }
}