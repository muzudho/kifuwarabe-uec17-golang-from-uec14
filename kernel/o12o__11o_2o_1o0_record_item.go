package kernel

import point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"

// RecordItem - 棋譜の一手分
type RecordItem struct {
	// 着手点
	placePlay point.Point

	// [O22o7o1o0] コウの位置
	ko point.Point
}

// NewRecordItem - 棋譜の一手分
func NewRecordItem() *RecordItem {
	var ri = new(RecordItem)
	return ri
}

// Clear - 空っぽにします
func (ri *RecordItem) Clear() {
	ri.placePlay = point.Point(0)
	ri.ko = point.Point(0)
}
