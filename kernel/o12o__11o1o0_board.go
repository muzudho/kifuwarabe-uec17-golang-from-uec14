package kernel

import (
	types1 "github.com/muzudho/kifuwarabe-uec17/kernel/types1"
	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
)

// Board - 盤
type Board struct {
	// ゲームルール
	gameRule GameRule
	// 盤座標
	coordinate types2.BoardCoordinate

	// 交点
	//
	// * 英語で交点は node かも知れないが、表計算でよく使われる cell の方を使う
	cells []types2.Stone
}

// NewBoard - 新規作成
func NewBoard(gameRule GameRule, boardWidht int, boardHeight int) *Board {
	var b = new(Board)

	// 設定ファイルから読込むので動的設定
	b.gameRule = gameRule

	// 枠の分、２つ増える
	var memoryBoardWidth = boardWidht + 2
	var memoryBoardHeight = boardHeight + 2

	b.coordinate = types2.BoardCoordinate{
		MemoryWidth:  memoryBoardWidth,
		MemoryHeight: memoryBoardHeight,
		// ４方向（東、北、西、南）の番地への相対インデックス
		Cell4Directions: [4]types1.Point{
			1,
			types1.Point(-memoryBoardWidth),
			-1,
			types1.Point(memoryBoardWidth),
		},
	}

	// 盤のサイズ指定と、盤面の初期化
	b.resize(boardWidht, boardHeight)

	return b
}

// GetGameRule - ゲームルール取得
func (b *Board) GetGameRule() *GameRule {
	return &b.gameRule
}

// SetGameRule - ゲームルール設定
func (b *Board) SetGameRule(gameRule *GameRule) {
	b.gameRule = *gameRule
}

// GetCoordinate - 盤座標取得
func (b *Board) GetCoordinate() *types2.BoardCoordinate {
	return &b.coordinate
}

// GetStoneAt - 指定座標の石を取得
func (b *Board) GetStoneAt(i types1.Point) types2.Stone {
	return b.cells[i]
}

// SetStoneAt - 指定座標の石を設定
func (b *Board) SetStoneAt(i types1.Point, s types2.Stone) {
	b.cells[i] = s
}

// GetColorAt - 指定座標の石の色を取得
func (b *Board) GetColorAt(i types1.Point) types1.Color {
	return b.cells[i].GetColor()
}

// IsEmpty - 指定の交点は空点か？
func (b *Board) IsSpaceAt(point types1.Point) bool {
	return b.GetStoneAt(point) == types2.Stone_Space
}

// サイズ変更
func (b *Board) resize(width int, height int) {
	b.coordinate.MemoryWidth = width + types2.BothSidesWallThickness
	b.coordinate.MemoryHeight = height + types2.BothSidesWallThickness
	b.cells = make([]types2.Stone, b.coordinate.GetMemoryArea())

	// ４方向（東、北、西、南）の番地への相対インデックス
	b.coordinate.Cell4Directions = [4]types1.Point{1, types1.Point(-b.coordinate.GetMemoryWidth()), -1, types1.Point(b.coordinate.GetMemoryWidth())}
}
