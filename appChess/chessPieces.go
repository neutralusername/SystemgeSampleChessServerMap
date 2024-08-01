package appChess

type Piece interface {
	isWhite() bool
	getLetter() string
}

type Knight struct {
	white bool
}

func (k *Knight) isWhite() bool {
	return k.white
}

func (k *Knight) getLetter() string {
	if k.white {
		return "N"
	}
	return "n"
}

type Rook struct {
	white    bool
	hasMoved bool
}

func (r *Rook) isWhite() bool {
	return r.white
}

func (r *Rook) getLetter() string {
	if r.white {
		return "R"
	}
	return "r"
}

type Bishop struct {
	white bool
}

func (b *Bishop) isWhite() bool {
	return b.white
}

func (b *Bishop) getLetter() string {
	if b.white {
		return "B"
	}
	return "b"
}

type Queen struct {
	white bool
}

func (q *Queen) isWhite() bool {
	return q.white
}

func (q *Queen) getLetter() string {
	if q.white {
		return "Q"
	}
	return "q"
}

type King struct {
	white    bool
	hasMoved bool
}

func (k *King) isWhite() bool {
	return k.white
}

func (k *King) getLetter() string {
	if k.white {
		return "K"
	}
	return "k"
}

type Pawn struct {
	white bool
}

func (p *Pawn) isWhite() bool {
	return p.white
}

func (p *Pawn) getLetter() string {
	if p.white {
		return "P"
	}
	return "p"
}
