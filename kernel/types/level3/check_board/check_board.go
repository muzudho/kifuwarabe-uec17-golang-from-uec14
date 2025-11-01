package check_board

import (
	// Level 2.1
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 4.1
	mark "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_technical_conceptual/sublevel_1/mark"

	// Level 2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_2/board_coordinate"
)

// CheckBoard - チェック盤
type CheckBoard struct {
	// 盤座標
	coordinate board_coordinate.BoardCoordinate

	// 長さが可変な盤
	//
	// * 英語で交点は node かも知れないが、表計算でよく使われる cell の方を使う
	cells []mark.Mark
}

// NewDirtyCheckBoard - 新規作成するが、初期化されていない
//
// * このメソッドを呼び出した後に Init 関数を呼び出してほしい
func NewDirtyCheckBoard() *CheckBoard {
	var cb = new(CheckBoard)

	cb.coordinate = board_coordinate.BoardCoordinate{}

	return cb
}

// Init - 初期化
func (cb *CheckBoard) Init(newBoardCoordinate board_coordinate.BoardCoordinate) {
	// 盤面のサイズが異なるなら、盤面を作り直す
	if cb.coordinate.MemoryWidth != newBoardCoordinate.MemoryWidth || cb.coordinate.MemoryHeight != newBoardCoordinate.MemoryHeight {
		cb.coordinate = newBoardCoordinate
		cb.cells = make([]mark.Mark, cb.coordinate.GetMemoryArea())
		return
	}

	// 盤面のクリアー
	for p := point.Point(0); p < point.Point(len(cb.cells)); p++ {
		cb.cells[p] = mark.Mark_BitAllZeros
	}
}

// GetAllBitsAt - 指定した交点の目印を取得
func (cb *CheckBoard) GetAllBitsAt(point point.Point) mark.Mark {
	return cb.cells[point]
}

// SetAllBitsAt - 指定した交点に目印を設定
func (cb *CheckBoard) SetAllBitsAt(point point.Point, mark mark.Mark) {
	cb.cells[point] = mark
}

// ClearAllBitsAt - フラグを消す
func (cb *CheckBoard) ClearAllBitsAt(point point.Point) {
	cb.cells[point] = mark.Mark(0)
}

// IsZeroAt - 指定した交点に目印は付いていないか？
func (cb *CheckBoard) IsZeroAt(point point.Point) bool {
	return cb.cells[point] == mark.Mark_BitAllZeros
}

// Overwrite - 上書き
func (cb *CheckBoard) Overwrite(point point.Point, mark mark.Mark) {
	cb.cells[point] |= mark
}

// Erase - 消す
func (cb *CheckBoard) Erase(point point.Point, mark mark.Mark) {
	cb.cells[point] &= ^mark // ^ はビット反転
}

// Contains - 含む
func (cb *CheckBoard) Contains(point point.Point, mark mark.Mark) bool {
	return cb.cells[point]&mark == mark
}
