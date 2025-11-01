package kernel

import (
	point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"

	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
)

// Mark - 目印
type Mark uint8

const (
	Mark_BitAllZeros Mark = 0b00000000
	Mark_BitStone    Mark = 0b00000001
	Mark_BitLiberty  Mark = 0b00000010
)

// CheckBoard - チェック盤
type CheckBoard struct {
	// 盤座標
	coordinate types2.BoardCoordinate

	// 長さが可変な盤
	//
	// * 英語で交点は node かも知れないが、表計算でよく使われる cell の方を使う
	cells []Mark
}

// NewDirtyCheckBoard - 新規作成するが、初期化されていない
//
// * このメソッドを呼び出した後に Init 関数を呼び出してほしい
func NewDirtyCheckBoard() *CheckBoard {
	var cb = new(CheckBoard)

	cb.coordinate = types2.BoardCoordinate{}

	return cb
}

// Init - 初期化
func (cb *CheckBoard) Init(newBoardCoordinate types2.BoardCoordinate) {
	// 盤面のサイズが異なるなら、盤面を作り直す
	if cb.coordinate.MemoryWidth != newBoardCoordinate.MemoryWidth || cb.coordinate.MemoryHeight != newBoardCoordinate.MemoryHeight {
		cb.coordinate = newBoardCoordinate
		cb.cells = make([]Mark, cb.coordinate.GetMemoryArea())
		return
	}

	// 盤面のクリアー
	for p := point.Point(0); p < point.Point(len(cb.cells)); p++ {
		cb.cells[p] = Mark_BitAllZeros
	}
}

// GetAllBitsAt - 指定した交点の目印を取得
func (cb *CheckBoard) GetAllBitsAt(point point.Point) Mark {
	return cb.cells[point]
}

// SetAllBitsAt - 指定した交点に目印を設定
func (cb *CheckBoard) SetAllBitsAt(point point.Point, mark Mark) {
	cb.cells[point] = mark
}

// ClearAllBitsAt - フラグを消す
func (cb *CheckBoard) ClearAllBitsAt(point point.Point) {
	cb.cells[point] = Mark(0)
}

// IsZeroAt - 指定した交点に目印は付いていないか？
func (cb *CheckBoard) IsZeroAt(point point.Point) bool {
	return cb.cells[point] == Mark_BitAllZeros
}

// Overwrite - 上書き
func (cb *CheckBoard) Overwrite(point point.Point, mark Mark) {
	cb.cells[point] |= mark
}

// Erase - 消す
func (cb *CheckBoard) Erase(point point.Point, mark Mark) {
	cb.cells[point] &= ^mark // ^ はビット反転
}

// Contains - 含む
func (cb *CheckBoard) Contains(point point.Point, mark Mark) bool {
	return cb.cells[point]&mark == mark
}
