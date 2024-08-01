package appChess

import (
	"SystemgeSampleChessServer/dto"

	"github.com/neutralusername/Systemge/Error"
)

func (game *ChessGame) move(move *dto.Move) (*dto.Move, error) {
	piece := game.board[move.FromRow][move.FromCol]
	if piece == nil {
		return nil, Error.New("no piece at from coordinates", nil)
	}
	if game.isWhiteTurn() != piece.isWhite() {
		return nil, Error.New("Cannot move opponent's piece", nil)
	}
	if err := game.isLegalMove(move.FromRow, move.FromCol, move.ToRow, move.ToCol); err != nil {
		return nil, Error.New("Illegal move", err)
	}
	notation := game.generateAlgebraicNotation(move.FromRow, move.FromCol, move.ToRow, move.ToCol)
	switch piece.(type) {
	case *King:
		if move.FromCol-move.ToCol == 2 {
			game.board[move.FromRow][move.FromCol-4], game.board[move.FromRow][move.FromCol-1] = game.board[move.FromRow][move.FromCol-1], game.board[move.FromRow][move.FromCol-4]
		} else if move.FromCol-move.ToCol == -2 {
			game.board[move.FromRow][move.FromCol+3], game.board[move.FromRow][move.FromCol+1] = game.board[move.FromRow][move.FromCol+1], game.board[move.FromRow][move.FromCol+3]
		}
		piece.(*King).hasMoved = true
	case *Pawn:
		if move.FromCol != move.ToCol && game.board[move.ToRow][move.ToCol] == nil {
			game.board[move.ToRow-1][move.ToCol] = nil
		}
		if move.ToRow == 0 || move.ToRow == 7 {
			game.board[move.FromRow][move.FromCol] = &Queen{white: piece.isWhite()}
		}
	case *Rook:
		piece.(*Rook).hasMoved = true
	}
	game.board[move.ToRow][move.ToCol] = game.board[move.FromRow][move.FromCol]
	game.board[move.FromRow][move.FromCol] = nil
	move.AlgebraicNotation = notation
	game.moves = append(game.moves, move)
	return move, nil
}

func (game *ChessGame) isLegalMove(fromRow, fromCol, toRow, toCol int) error {
	if fromRow < 0 || fromRow > 7 || fromCol < 0 || fromCol > 7 || toRow < 0 || toRow > 7 || toCol < 0 || toCol > 7 {
		return Error.New("coordinates out of bounds", nil)
	}
	fromPece := game.board[fromRow][fromCol]
	toPiece := game.board[toRow][toCol]
	if toPiece != nil && toPiece.isWhite() == fromPece.isWhite() {
		return Error.New("cannot take own piece", nil)
	}
	switch fromPece.(type) {
	case *Rook:
		if err := game.isValidRookMove(fromRow, fromCol, toRow, toCol); err != nil {
			return Error.New("invalid rook move", err)
		}
	case *Bishop:
		if err := game.isValidBishopMove(fromRow, fromCol, toRow, toCol); err != nil {
			return Error.New("invalid bishop move", err)
		}
	case *Queen:
		if err := game.isValidRookMove(fromRow, fromCol, toRow, toCol); err != nil {
			if err := game.isValidBishopMove(fromRow, fromCol, toRow, toCol); err != nil {
				return Error.New("invalid queen move", err)
			}
		}
	case *King:
		if err := game.isValidKingMove(fromRow, fromCol, toRow, toCol); err != nil {
			if err := game.isValidCastleMove(fromRow, fromCol, toRow, toCol); err != nil {
				return Error.New("invalid king move", err)
			}
		}
	case *Pawn:
		if err := game.isValidPawnMove(fromRow, fromCol, toRow, toCol); err != nil {
			return Error.New("invalid pawn move", err)
		}
	case *Knight:
		if err := game.isValidKnightMove(fromRow, fromCol, toRow, toCol); err != nil {
			return Error.New("invalid knight move", err)
		}
	}
	if game.isInCheckAfterMove(fromRow, fromCol, toRow, toCol) {
		return Error.New("cannot move into check", nil)
	}
	return nil
}

func (game *ChessGame) isInCheckAfterMove(fromRow, fromCol, toRow, toCol int) bool {
	kingRow, kingCol := game.getKingCoordinates(game.isWhiteTurn())
	if kingRow == -1 || kingCol == -1 {
		return false
	}
	if fromRow == kingRow && fromCol == kingCol {
		kingRow, kingCol = toRow, toCol
	}
	kingPiece := game.board[kingRow][kingCol]
	game.board[kingRow][kingCol] = nil
	game.board[toRow][toCol] = game.board[fromRow][fromCol]
	game.board[fromRow][fromCol] = nil
	defer func() {
		game.board[fromRow][fromCol] = game.board[toRow][toCol]
		game.board[toRow][toCol] = nil
		game.board[kingRow][kingCol] = kingPiece
	}()
	for i, row := range game.board {
		for j, piece := range row {
			if piece != nil && piece.isWhite() != game.isWhiteTurn() {
				if err := game.isLegalMove(i, j, kingRow, kingCol); err == nil {
					return true
				}
			}
		}
	}
	return false
}

func (game *ChessGame) getKingCoordinates(isWhite bool) (int, int) {
	for i, row := range game.board {
		for j, piece := range row {
			if king, ok := piece.(*King); ok {
				if king.isWhite() == isWhite {
					return i, j
				}
			}
		}
	}
	return -1, -1
}

func (game *ChessGame) isValidKnightMove(fromRow, fromCol, toRow, toCol int) error {
	if (fromRow-toRow != 2 && fromRow-toRow != -2) || (fromCol-toCol != 1 && fromCol-toCol != -1) {
		if (fromRow-toRow != 1 && fromRow-toRow != -1) || (fromCol-toCol != 2 && fromCol-toCol != -2) {
			return Error.New("knight can only move in L shape", nil)
		}
	}
	return nil
}

func (game *ChessGame) isValidPawnMove(fromRow, fromCol, toRow, toCol int) error {
	fromPiece := game.board[fromRow][fromCol].(*Pawn)
	toPiece := game.board[toRow][toCol]

	if fromPiece.isWhite() {
		if fromCol == toCol {
			if fromRow-toRow == -1 {
				return nil
			} else if fromRow-toRow == -2 {
				if fromRow == 1 {
					return nil
				} else {
					return Error.New("pawn can only move two squares on first move", nil)
				}
			} else if fromRow-toRow < -2 {
				return Error.New("pawn cannot move more than two squares", nil)
			} else if fromRow-toRow == 0 {
				return Error.New("pawn cannot move horizontally", nil)
			} else if fromRow-toRow > 0 {
				return Error.New("pawn cannot move backwards", nil)
			}
		}
		if (fromCol-toCol != 1 && fromCol-toCol != -1) || fromRow-toRow != -1 {
			return Error.New("pawn can only move one square diagonally to take a piece", nil)
		}
		if toPiece != nil {
			return nil
		}
		if fromRow != 4 {
			return Error.New("can only en passant from fifth rank", nil)
		}
		lastMove := game.moves[len(game.moves)-1]
		lastPiece := game.board[lastMove.ToRow][lastMove.ToCol]
		if _, ok := lastPiece.(*Pawn); ok {
			if lastMove.ToRow-lastMove.FromRow == -2 && lastMove.ToCol == toCol {
				return nil
			}
		}
		return Error.New("can only en passant immediately after opponent's pawn moves two squares", nil)
	} else {
		if fromCol == toCol {
			if fromRow-toRow == 1 {
				return nil
			} else if fromRow-toRow == 2 {
				if fromRow == 6 {
					return nil
				} else {
					return Error.New("pawn can only move two squares on first move", nil)
				}
			} else if fromRow-toRow > 2 {
				return Error.New("pawn cannot move more than two squares", nil)
			} else if fromRow-toRow == 0 {
				return Error.New("pawn cannot move horizontally", nil)
			} else if fromRow-toRow < 0 {
				return Error.New("pawn cannot move backwards", nil)
			}
		}
		if (fromCol-toCol != 1 && fromCol-toCol != -1) || fromRow-toRow != 1 {
			return Error.New("pawn can only move one square diagonally to take a piece", nil)
		}
		if toPiece != nil {
			return nil
		}
		if fromRow != 3 {
			return Error.New("can only en passant from fourth rank", nil)
		}
		lastMove := game.moves[len(game.moves)-1]
		lastPiece := game.board[lastMove.ToRow][lastMove.ToCol]
		if _, ok := lastPiece.(*Pawn); ok {
			if lastMove.ToRow-lastMove.FromRow == 2 && lastMove.ToCol == toCol {
				return nil
			}
		}
		return Error.New("can only en passant immediately after opponent's pawn moves two squares", nil)
	}
}

func (game *ChessGame) isValidKingMove(fromRow, fromCol, toRow, toCol int) error {
	if fromRow-toRow > 1 || fromRow-toRow < -1 || fromCol-toCol > 1 || fromCol-toCol < -1 {
		return Error.New("king can only move one square in any direction", nil)
	}
	return nil
}

func (game *ChessGame) isValidCastleMove(fromRow, fromCol, toRow, toCol int) error {
	king := game.board[fromRow][fromCol].(*King)
	if king.hasMoved {
		return Error.New("king has already moved", nil)
	}
	if fromRow != toRow {
		return Error.New("king can only castle horizontally", nil)
	}
	if fromCol-toCol == 2 {
		rook := game.board[fromRow][0].(*Rook)
		if rook.hasMoved {
			return Error.New("rook has already moved", nil)
		}
		for i := 1; i < 4; i++ {
			if game.board[fromRow][i] != nil {
				return Error.New("cannot castle through pieces", nil)
			}
		}
		for i := 3; i < 5; i++ {
			if game.isInCheckAfterMove(fromRow, fromCol, fromRow, i) {
				return Error.New("cannot castle through check", nil)
			}
		}
	}
	if fromCol-toCol == -2 {
		rook := game.board[fromRow][7].(*Rook)
		if rook.hasMoved {
			return Error.New("rook has already moved", nil)
		}
		for i := 5; i < 7; i++ {
			if game.board[fromRow][i] != nil {
				return Error.New("cannot castle through pieces", nil)
			}
		}
		for i := 4; i < 6; i++ {
			if game.isInCheckAfterMove(fromRow, fromCol, fromRow, i) {
				return Error.New("cannot castle through check", nil)
			}
		}
	}
	return nil
}

func (game *ChessGame) isValidBishopMove(fromRow, fromCol, toRow, toCol int) error {
	if fromRow-toRow != fromCol-toCol && fromRow-toRow != toCol-fromCol {
		return Error.New("bishop can only move diagonally", nil)
	}
	if fromRow < toRow {
		if fromCol < toCol {
			for i, j := fromRow+1, fromCol+1; i < toRow; i, j = i+1, j+1 {
				if game.board[i][j] != nil {
					return Error.New("bishop cannot jump over pieces", nil)
				}
			}
		} else {
			for i, j := fromRow+1, fromCol-1; i < toRow; i, j = i+1, j-1 {
				if game.board[i][j] != nil {
					return Error.New("bishop cannot jump over pieces", nil)
				}
			}
		}
	} else {
		if fromCol < toCol {
			for i, j := fromRow-1, fromCol+1; i > toRow; i, j = i-1, j+1 {
				if game.board[i][j] != nil {
					return Error.New("bishop cannot jump over pieces", nil)
				}
			}
		} else {
			for i, j := fromRow-1, fromCol-1; i > toRow; i, j = i-1, j-1 {
				if game.board[i][j] != nil {
					return Error.New("bishop cannot jump over pieces", nil)
				}
			}
		}
	}
	return nil

}

func (game *ChessGame) isValidRookMove(fromRow, fromCol, toRow, toCol int) error {
	if fromRow != toRow && fromCol != toCol {
		return Error.New("rook can only move in a straight line", nil)
	}
	if fromRow == toRow {
		if fromCol < toCol {
			for i := fromCol + 1; i < toCol; i++ {
				if game.board[fromRow][i] != nil {
					return Error.New("rook cannot jump over pieces", nil)
				}
			}
		} else {
			for i := fromCol - 1; i > toCol; i-- {
				if game.board[fromRow][i] != nil {
					return Error.New("rook cannot jump over pieces", nil)
				}
			}
		}
	} else {
		if fromRow < toRow {
			for i := fromRow + 1; i < toRow; i++ {
				if game.board[i][fromCol] != nil {
					return Error.New("rook cannot jump over pieces", nil)
				}
			}
		} else {
			for i := fromRow - 1; i > toRow; i-- {
				if game.board[i][fromCol] != nil {
					return Error.New("rook cannot jump over pieces", nil)
				}
			}
		}
	}
	return nil
}
