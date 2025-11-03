package record

import (
	"math"
	"strconv"

	// Entities
	moves_num "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/moves_num"
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"
	stone "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/stone"

	// Level 2.1
	geta "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_1/geta"

	// Level 6.1
	record_item "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_6_database/sublevel_1/record_item"
)

// Record - 棋譜
type Record struct {
	// 先行
	playFirst stone.Stone

	// 何手目。基数
	MovesNum1 moves_num.MovesNum

	// 手毎
	items []*record_item.RecordItem
}

// NewRecord - 新規作成
//
// * maxPositionNumber - 手数上限。配列サイズ決定のための判断材料
// * memoryBoardArea - メモリー盤サイズ。配列サイズ決定のための判断材料
func NewRecord(movesNum1 moves_num.MovesNum, memoryBoardArea int, playFirst stone.Stone) *Record {
	var r = new(Record)
	r.playFirst = playFirst

	// 動的に長さがきまる配列を生成、その内容をインスタンスで埋めます
	// 例えば、0手目が初期局面として、 400 手目まであるとすると、要素数は401要る。だから 1 足す
	// しかし、プレイアウトでは終局まで打ちたいので、多めにとっておきたいのでは。盤サイズより適当に18倍（>2πe）取る
	var positionLength = int(math.Max(float64(movesNum1+1), float64(memoryBoardArea*18)))
	r.items = make([]*record_item.RecordItem, positionLength)

	for i := moves_num.MovesNum(0); i < moves_num.MovesNum(positionLength); i++ {
		r.items[i] = record_item.NewRecordItem()
	}

	return r
}

// GetMaxPosNthFigure - 手数（序数）の最大値の桁数
func (r *Record) GetMaxPosNthFigure() int {
	var nth = r.GetMaxPosNth()
	var nthText = strconv.Itoa(nth)
	return len(nthText)
}

// GetMaxPosNth - 手数（序数）の最大値
func (r *Record) GetMaxPosNth() int {
	return len(r.items) + geta.Geta
}

// GetMovesNum - 何手目。基数
func (r *Record) GetMovesNum() moves_num.MovesNum {
	return r.MovesNum1
}

// Push - 末尾に追加
func (r *Record) Push(placePlay point.Point,
	// [O22o7o1o0] コウの位置
	ko point.Point) {

	var item = r.items[r.MovesNum1]
	item.PlacePlay = placePlay

	// [O22o7o1o0] コウの位置
	item.Ko = ko

	r.MovesNum1++
}

// RemoveTail - 末尾を削除
func (r *Record) RemoveTail(placePlay point.Point) {
	r.MovesNum1--
	r.items[r.MovesNum1].Clear()
}

// ForeachItem - 各要素
func (r *Record) ForeachItem(setItem func(moves_num.MovesNum, *record_item.RecordItem)) {
	for i := moves_num.MovesNum(0); i < r.MovesNum1; i++ {
		setItem(i, r.items[i])
	}
}

// IsKo - コウか？
func (r *Record) IsKo(placePlay point.Point) bool {
	// [O22o7o1o0] コウの判定
	// 2手前に着手して石をぴったり１つ打ち上げたとき、その着手点はコウだ
	var positionNumber = r.GetMovesNum()
	if 2 <= positionNumber {
		var item = r.items[positionNumber-2]
		return item.Ko == placePlay
	}

	return false
}
