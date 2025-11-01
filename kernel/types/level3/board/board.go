package board

import (
	// Level 2.1
	color "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/color"
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_2/board_coordinate"
	game_rule_settings "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/game_rule_settings"
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/stone"
)

// Board - 盤
type Board struct {
	// gameRuleSettings - 対局ルール設定
	gameRuleSettings game_rule_settings.GameRuleSettings

	// Coordinate - 盤座標
	Coordinate board_coordinate.BoardCoordinate

	// Cells - 交点
	//
	// * 英語で交点は node かも知れないが、表計算でよく使われる cell の方を使う
	Cells []stone.Stone
}

// NewBoard - 新規作成
func NewBoard(gameRuleSettings game_rule_settings.GameRuleSettings, boardWidht int, boardHeight int) *Board {
	var b = new(Board)

	// 設定ファイルから読込むので動的設定
	b.gameRuleSettings = gameRuleSettings

	// 枠の分、２つ増える
	var memoryBoardWidth = boardWidht + 2
	var memoryBoardHeight = boardHeight + 2

	b.Coordinate = board_coordinate.BoardCoordinate{
		MemoryWidth:  memoryBoardWidth,
		MemoryHeight: memoryBoardHeight,
		// ４方向（東、北、西、南）の番地への相対インデックス
		Cell4Directions: [4]point.Point{
			1,
			point.Point(-memoryBoardWidth),
			-1,
			point.Point(memoryBoardWidth),
		},
	}

	// 盤のサイズ指定と、盤面の初期化
	b.resize(boardWidht, boardHeight)

	return b
}

// GetGameRule - ゲームルール取得
func (b *Board) GetGameRule() *game_rule_settings.GameRuleSettings {
	return &b.gameRuleSettings
}

// SetGameRule - ゲームルール設定
func (b *Board) SetGameRule(gameRuleSettings *game_rule_settings.GameRuleSettings) {
	b.gameRuleSettings = *gameRuleSettings
}

// GetCoordinate - 盤座標取得
func (b *Board) GetCoordinate() *board_coordinate.BoardCoordinate {
	return &b.Coordinate
}

// GetStoneAt - 指定座標の石を取得
func (b *Board) GetStoneAt(i point.Point) stone.Stone {
	return b.Cells[i]
}

// SetStoneAt - 指定座標の石を設定
func (b *Board) SetStoneAt(i point.Point, s stone.Stone) {
	b.Cells[i] = s
}

// GetColorAt - 指定座標の石の色を取得
func (b *Board) GetColorAt(i point.Point) color.Color {
	return b.Cells[i].GetColor()
}

// IsEmpty - 指定の交点は空点か？
func (b *Board) IsSpaceAt(point point.Point) bool {
	return b.GetStoneAt(point) == stone.Stone_Space
}

// サイズ変更
func (b *Board) resize(width int, height int) {
	b.Coordinate.MemoryWidth = width + board_coordinate.BothSidesWallThickness
	b.Coordinate.MemoryHeight = height + board_coordinate.BothSidesWallThickness
	b.Cells = make([]stone.Stone, b.Coordinate.GetMemoryArea())

	// ４方向（東、北、西、南）の番地への相対インデックス
	b.Coordinate.Cell4Directions = [4]point.Point{1, point.Point(-b.Coordinate.GetMemoryWidth()), -1, point.Point(b.Coordinate.GetMemoryWidth())}
}

// Init - 盤面初期化
func (b *Board) Init(width int, height int) {
	// 盤面のサイズが異なるなら、盤面を作り直す
	if b.Coordinate.MemoryWidth != width || b.Coordinate.MemoryHeight != height {
		b.resize(width, height)
	}

	// 枠の上辺、下辺を引く
	{
		var y = 0
		var y2 = b.Coordinate.MemoryHeight - 1
		for x := 0; x < b.Coordinate.MemoryWidth; x++ {
			var i = b.Coordinate.GetPointFromXy(x, y)
			b.Cells[i] = stone.Stone_Wall

			i = b.Coordinate.GetPointFromXy(x, y2)
			b.Cells[i] = stone.Stone_Wall
		}
	}
	// 枠の左辺、右辺を引く
	{
		var x = 0
		var x2 = b.Coordinate.MemoryWidth - 1
		for y := 1; y < b.Coordinate.MemoryHeight-1; y++ {
			var i = b.Coordinate.GetPointFromXy(x, y)
			b.Cells[i] = stone.Stone_Wall

			i = b.Coordinate.GetPointFromXy(x2, y)
			b.Cells[i] = stone.Stone_Wall
		}
	}
	// 枠の内側を空点で埋める
	{
		var height = b.Coordinate.GetHeight()
		var width = b.Coordinate.GetWidth()
		for y := 1; y < height; y++ {
			for x := 1; x < width; x++ {
				var i = b.Coordinate.GetPointFromXy(x, y)
				b.Cells[i] = stone.Stone_Space
			}
		}
	}
}

// ForeachNeumannNeighborhood - [O13o__10o0] 隣接する４方向の定義
func (b *Board) ForeachNeumannNeighborhood(here point.Point, setAdjacent func(board_coordinate.Cell_4Directions, point.Point)) {
	// 東、北、西、南
	for dir := board_coordinate.Cell_4Directions(0); dir < 4; dir++ {
		var p = here + b.Coordinate.Cell4Directions[dir] // 隣接する交点

		// 範囲外チェック
		if p < 0 || b.Coordinate.GetMemoryArea() <= int(p) {
			continue
		}

		// 枠チェック
		if b.GetStoneAt(p) == stone.Stone_Wall {
			continue
		}

		setAdjacent(dir, p)
	}
}

// IsMasonry - 石の上に石を置こうとしたか？
func (b *Board) IsMasonry(point point.Point) bool {
	// 空点以外に石を置こうとしたら、石の上に石を置いた扱いにする
	return !b.IsSpaceAt(point)
}
