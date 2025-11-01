// BOF [O12o__11o_2o0]

package kernel

import (
	"math"
	"strconv"

	point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"
	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
)

// Record - 棋譜
type Record struct {
	// 先行
	playFirst types2.Stone

	// 何手目。基数
	positionNumber PositionNumberInt

	// 手毎
	items []*RecordItem
}

// NewRecord - 新規作成
//
// * maxPositionNumber - 手数上限。配列サイズ決定のための判断材料
// * memoryBoardArea - メモリー盤サイズ。配列サイズ決定のための判断材料
func NewRecord(maxPositionNumber PositionNumberInt, memoryBoardArea int, playFirst types2.Stone) *Record {
	var r = new(Record)
	r.playFirst = playFirst

	// 動的に長さがきまる配列を生成、その内容をインスタンスで埋めます
	// 例えば、0手目が初期局面として、 400 手目まであるとすると、要素数は401要る。だから 1 足す
	// しかし、プレイアウトでは終局まで打ちたいので、多めにとっておきたいのでは。盤サイズより適当に18倍（>2πe）取る
	var positionLength = int(math.Max(float64(maxPositionNumber+1), float64(memoryBoardArea*18)))
	r.items = make([]*RecordItem, positionLength)

	for i := PositionNumberInt(0); i < PositionNumberInt(positionLength); i++ {
		r.items[i] = NewRecordItem()
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
	return len(r.items) + geta
}

// GetPositionNumber - 何手目。基数
func (r *Record) GetPositionNumber() PositionNumberInt {
	return r.positionNumber
}

// Push - 末尾に追加
func (r *Record) Push(placePlay point.Point,
	// [O22o7o1o0] コウの位置
	ko point.Point) {

	var item = r.items[r.positionNumber]
	item.placePlay = placePlay

	// [O22o7o1o0] コウの位置
	item.ko = ko

	r.positionNumber++
}

// RemoveTail - 末尾を削除
func (r *Record) RemoveTail(placePlay point.Point) {
	r.positionNumber--
	r.items[r.positionNumber].Clear()
}

// ForeachItem - 各要素
func (r *Record) ForeachItem(setItem func(PositionNumberInt, *RecordItem)) {
	for i := PositionNumberInt(0); i < r.positionNumber; i++ {
		setItem(i, r.items[i])
	}
}

// IsKo - コウか？
func (r *Record) IsKo(placePlay point.Point) bool {
	// [O22o7o1o0] コウの判定
	// 2手前に着手して石をぴったり１つ打ち上げたとき、その着手点はコウだ
	var positionNumber = r.GetPositionNumber()
	if 2 <= positionNumber {
		var item = r.items[positionNumber-2]
		return item.ko == placePlay
	}

	return false
}
