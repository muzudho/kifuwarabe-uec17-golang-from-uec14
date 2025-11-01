package kernel

import types1 "github.com/muzudho/kifuwarabe-uec17/kernel/types1"

// RecordItem - 棋譜の一手分
type RecordItem struct {
	// 着手点
	placePlay types1.Point

	// [O22o7o1o0] コウの位置
	ko types1.Point
}

// NewRecordItem - 棋譜の一手分
func NewRecordItem() *RecordItem {
	var ri = new(RecordItem)
	return ri
}

// Clear - 空っぽにします
func (ri *RecordItem) Clear() {
	ri.placePlay = types1.Point(0)
	ri.ko = types1.Point(0)
}
