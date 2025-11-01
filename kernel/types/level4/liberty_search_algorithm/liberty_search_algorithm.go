package liberty_search_algorithm

import (
	// Level 2.1
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2.2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_2/board_coordinate"

	// Level 3.1
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/level_3_physical/sublevel_1/stone"

	// Level 4.1
	check_board "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_game_rule/sublevel_1/check_board"
	rentype "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_game_rule/sublevel_1/ren"

	// Level 4.2
	board "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_game_rule/sublevel_2/board"

	// Level 5.1
	mark "github.com/muzudho/kifuwarabe-uec17/kernel/level_5_game_technic/sublevel_1/mark"
)

// LibertySearchAlgorithm - 呼吸点探索アルゴリズム
type LibertySearchAlgorithm struct {
	// 盤
	board *board.Board
	// チェック盤
	checkBoard *check_board.CheckBoard
	// foundRen - 呼吸点の探索時に使います
	foundRen *rentype.Ren
}

// NewLibertySearchAlgorithm - 新規作成
func NewLibertySearchAlgorithm(board *board.Board, checkBoard1 *check_board.CheckBoard) *LibertySearchAlgorithm {
	var ls = new(LibertySearchAlgorithm)

	ls.board = board
	ls.checkBoard = checkBoard1

	return ls
}

// 連の検索
//
// Returns
// -------
// - *Ren is ren or nil
// - bool is found
func (ls *LibertySearchAlgorithm) FindRen(arbitraryPoint point.Point) (*rentype.Ren, bool) {
	// 探索済みならスキップ
	if ls.checkBoard.Contains(arbitraryPoint, mark.Mark_BitStone) {
		return nil, false
	}

	// 連の初期化
	ls.foundRen = rentype.NewRen(ls.board.GetStoneAt(arbitraryPoint))

	if ls.foundRen.Stone == stone.Stone_Space {
		ls.searchSpaceRen(arbitraryPoint)
	} else {
		ls.searchStoneRenRecursive(arbitraryPoint)

		// チェックボードの「呼吸点」チェックのみクリアー
		var eachPoint = func(point point.Point) {
			ls.checkBoard.Erase(point, mark.Mark_BitLiberty)
		}
		ls.board.Coordinate.ForeachCellWithoutWall(eachPoint)
	}

	return ls.foundRen, true
}

// 石の連の探索
//
// * 再帰関数
func (ls *LibertySearchAlgorithm) searchStoneRenRecursive(here point.Point) {

	// 石のチェック
	ls.checkBoard.Overwrite(here, mark.Mark_BitStone)

	ls.foundRen.AddLocation(here)

	// 隣接する交点毎に
	var eachAdjacent = func(dir board_coordinate.Cell_4Directions, p point.Point) {

		var stone1 = ls.board.GetStoneAt(p) // 石の色
		switch stone1 {

		case stone.Stone_Space: // 空点
			if !ls.checkBoard.Contains(p, mark.Mark_BitLiberty) { // まだチェックしていない呼吸点なら
				ls.checkBoard.Overwrite(p, mark.Mark_BitLiberty)
				ls.foundRen.LibertyLocations = append(ls.foundRen.LibertyLocations, p) // 呼吸点を追加
			}

			return // あとの処理をスキップ

		case stone.Stone_Wall: // 枠
			return // あとの処理をスキップ
		}

		// 探索済みの石ならスキップ
		if ls.checkBoard.Contains(p, mark.Mark_BitStone) {
			return
		}

		var color = stone1.GetColor()
		// 隣接する色、追加
		ls.foundRen.AdjacentColor = ls.foundRen.AdjacentColor.GetAdded(color)

		if stone1 == ls.foundRen.Stone { // 同じ石
			ls.searchStoneRenRecursive(p) // 再帰
		}
	}

	// 隣接する４方向
	ls.board.ForeachNeumannNeighborhood(here, eachAdjacent)
}

// 空点の連の探索
// - 再帰関数
func (ls *LibertySearchAlgorithm) searchSpaceRen(here point.Point) {
	ls.checkBoard.Overwrite(here, mark.Mark_BitStone)
	ls.foundRen.AddLocation(here)

	var eachAdjacent = func(dir board_coordinate.Cell_4Directions, p point.Point) {
		// 探索済みならスキップ
		if ls.checkBoard.Contains(p, mark.Mark_BitStone) {
			return
		}

		var stone1 = ls.board.GetStoneAt(p)
		if stone1 != stone.Stone_Space { // 空点でなければスキップ
			return
		}

		var color = stone1.GetColor()
		// 隣接する色、追加
		ls.foundRen.AdjacentColor = ls.foundRen.AdjacentColor.GetAdded(color)
		ls.searchSpaceRen(p) // 再帰
	}

	// 隣接する４方向
	ls.board.ForeachNeumannNeighborhood(here, eachAdjacent)
}
