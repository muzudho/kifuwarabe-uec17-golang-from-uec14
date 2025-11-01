package kernel

import (
	types1 "github.com/muzudho/kifuwarabe-uec17/kernel/types1"
	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
)

// GetLiberty - 呼吸点の数え上げ。連の数え上げ。
// `GetOneRen` とでもいう名前の方がふさわしいが、慣習に合わせた関数名にした
//
// Parameters
// ----------
// * `arbitraryPoint` - 連に含まれる任意の一点
//
// Returns
// -------
// - *Ren is ren or nil
// - bool is found
func (k *Kernel) GetLiberty(arbitraryPoint types1.Point) (*Ren, bool) {
	// チェックボードの初期化
	k.Position.CheckBoard.Init(k.Position.Board.coordinate)

	var libertySearchAlgorithm = NewLibertySearchAlgorithm(k.Position.Board, k.Position.CheckBoard)

	return libertySearchAlgorithm.findRen(arbitraryPoint)
}

// LibertySearchAlgorithm - 呼吸点探索アルゴリズム
type LibertySearchAlgorithm struct {
	// 盤
	board *Board
	// チェック盤
	checkBoard *CheckBoard
	// foundRen - 呼吸点の探索時に使います
	foundRen *Ren
}

// NewLibertySearchAlgorithm - 新規作成
func NewLibertySearchAlgorithm(board *Board, checkBoard *CheckBoard) *LibertySearchAlgorithm {
	var ls = new(LibertySearchAlgorithm)

	ls.board = board
	ls.checkBoard = checkBoard

	return ls
}

// 連の検索
//
// Returns
// -------
// - *Ren is ren or nil
// - bool is found
func (ls *LibertySearchAlgorithm) findRen(arbitraryPoint types1.Point) (*Ren, bool) {
	// 探索済みならスキップ
	if ls.checkBoard.Contains(arbitraryPoint, Mark_BitStone) {
		return nil, false
	}

	// 連の初期化
	ls.foundRen = NewRen(ls.board.GetStoneAt(arbitraryPoint))

	if ls.foundRen.stone == types2.Stone_Space {
		ls.searchSpaceRen(arbitraryPoint)
	} else {
		ls.searchStoneRenRecursive(arbitraryPoint)

		// チェックボードの「呼吸点」チェックのみクリアー
		var eachPoint = func(point types1.Point) {
			ls.checkBoard.Erase(point, Mark_BitLiberty)
		}
		ls.board.coordinate.ForeachCellWithoutWall(eachPoint)
	}

	return ls.foundRen, true
}

// 石の連の探索
//
// * 再帰関数
func (ls *LibertySearchAlgorithm) searchStoneRenRecursive(here types1.Point) {

	// 石のチェック
	ls.checkBoard.Overwrite(here, Mark_BitStone)

	ls.foundRen.AddLocation(here)

	// 隣接する交点毎に
	var eachAdjacent = func(dir Cell_4Directions, p types1.Point) {

		var stone = ls.board.GetStoneAt(p) // 石の色
		switch stone {

		case types2.Stone_Space: // 空点
			if !ls.checkBoard.Contains(p, Mark_BitLiberty) { // まだチェックしていない呼吸点なら
				ls.checkBoard.Overwrite(p, Mark_BitLiberty)
				ls.foundRen.libertyLocations = append(ls.foundRen.libertyLocations, p) // 呼吸点を追加
			}

			return // あとの処理をスキップ

		case types2.Stone_Wall: // 枠
			return // あとの処理をスキップ
		}

		// 探索済みの石ならスキップ
		if ls.checkBoard.Contains(p, Mark_BitStone) {
			return
		}

		var color = stone.GetColor()
		// 隣接する色、追加
		ls.foundRen.adjacentColor = ls.foundRen.adjacentColor.GetAdded(color)

		if stone == ls.foundRen.stone { // 同じ石
			ls.searchStoneRenRecursive(p) // 再帰
		}
	}

	// 隣接する４方向
	ls.board.ForeachNeumannNeighborhood(here, eachAdjacent)
}

// 空点の連の探索
// - 再帰関数
func (ls *LibertySearchAlgorithm) searchSpaceRen(here types1.Point) {
	ls.checkBoard.Overwrite(here, Mark_BitStone)
	ls.foundRen.AddLocation(here)

	var eachAdjacent = func(dir Cell_4Directions, p types1.Point) {
		// 探索済みならスキップ
		if ls.checkBoard.Contains(p, Mark_BitStone) {
			return
		}

		var stone = ls.board.GetStoneAt(p)
		if stone != types2.Stone_Space { // 空点でなければスキップ
			return
		}

		var color = stone.GetColor()
		// 隣接する色、追加
		ls.foundRen.adjacentColor = ls.foundRen.adjacentColor.GetAdded(color)
		ls.searchSpaceRen(p) // 再帰
	}

	// 隣接する４方向
	ls.board.ForeachNeumannNeighborhood(here, eachAdjacent)
}
