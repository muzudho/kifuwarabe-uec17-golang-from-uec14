package kernel

import types1 "github.com/muzudho/kifuwarabe-uec17/kernel/types1"

// IsMasonry - 石の上に石を置こうとしたか？
func (b *Board) IsMasonry(point types1.Point) bool {
	// 空点以外に石を置こうとしたら、石の上に石を置いた扱いにする
	return !b.IsSpaceAt(point)
}
