package level_31_controller

import (
	// Level 2.1
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 4.1
	rentype "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_game_rule/sublevel_1/ren"

	// Level 7.1
	liberty_search_algorithm "github.com/muzudho/kifuwarabe-uec17/kernel/level_7_misc/sublevel_1/liberty_search_algorithm"
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
func (k *Kernel) GetLiberty(arbitraryPoint point.Point) (*rentype.Ren, bool) {
	// チェックボードの初期化
	k.Position.CheckBoard.Init(k.Position.Board.Coordinate)

	var libertySearchAlgorithm = liberty_search_algorithm.NewLibertySearchAlgorithm(k.Position.Board, k.Position.CheckBoard)

	return libertySearchAlgorithm.FindRen(arbitraryPoint)
}
