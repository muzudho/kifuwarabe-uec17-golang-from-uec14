package level_31_controller

import (
	// Entities
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"
	ren "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_2/ren"

	// Level 7.1
	liberty_search_algorithm "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_7_misc/sublevel_1/liberty_search_algorithm"
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
func (kernel1 *Kernel) GetLiberty(arbitraryPoint point.Point) (*ren.Ren, bool) {
	// チェックボードの初期化
	kernel1.Position.CheckBoard.Init(kernel1.Position.Board.Coordinate)

	var libertySearchAlgorithm = liberty_search_algorithm.NewLibertySearchAlgorithm(kernel1.Position.Board, kernel1.Position.CheckBoard)

	return libertySearchAlgorithm.FindRen(arbitraryPoint)
}
