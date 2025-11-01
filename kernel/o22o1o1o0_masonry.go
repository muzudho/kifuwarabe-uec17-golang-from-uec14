package kernel

import point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"

// IsMasonry - 石の上に石を置こうとしたか？
func (b *Board) IsMasonry(point point.Point) bool {
	// 空点以外に石を置こうとしたら、石の上に石を置いた扱いにする
	return !b.IsSpaceAt(point)
}
