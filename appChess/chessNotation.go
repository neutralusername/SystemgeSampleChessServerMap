package appChess

import "github.com/neutralusername/Systemge/Helpers"

func (game *ChessGame) generateAlgebraicNotation(fromRow, fromCol, toRow, toCol int) string {
	notation := ""
	piece := game.board[fromRow][fromCol]
	switch piece.(type) {
	case *King:
		if fromCol-toCol == 2 {
			notation = "O-O-O"
		} else if fromCol-toCol == -2 {
			notation = "O-O"
		} else {
			notation = "K"
		}
		return notation
	case *Pawn:
		if fromCol != toCol && game.board[toRow][toCol] == nil {
			notation = game.getColumnLetter(fromCol) + "x" + game.getColumnLetter(toCol) + game.getRowNumber(toRow)
		} else {
			notation = game.getColumnLetter(toCol) + game.getRowNumber(toRow)
		}
		if toRow == 0 || toRow == 7 {
			notation += "=Q"
		}
		return notation
	}
	notation = piece.getLetter()
	if game.board[toRow][toCol] != nil {
		notation += "x"
	}
	notation += game.getColumnLetter(toCol) + game.getRowNumber(toRow)
	return notation
}

func (game *ChessGame) getColumnLetter(col int) string {
	switch col {
	case 0:
		return "a"
	case 1:
		return "b"
	case 2:
		return "c"
	case 3:
		return "d"
	case 4:
		return "e"
	case 5:
		return "f"
	case 6:
		return "g"
	case 7:
		return "h"
	}
	return ""
}

func (game *ChessGame) getRowNumber(row int) string {
	return Helpers.IntToString(8 - row)
}
