// BOF [O19o0]

package level_31_controller

import (
	"fmt"
	"math"
	"strings"

	// Section 1.1.1
	logger "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_7_presenter/chapter_1_i_o/section_1/logger"
	i_text_i_o "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/interfaces/part_1_facility/chapter_1_i_o/section_1/i_text_i_o"

	// Section 1.1.2

	// Level 2.1
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2.2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_2/board_coordinate"

	// Level 3.1
	stone "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_3_physical/sublevel_1/stone"

	// Level 4.1
	rentype "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_4_game_rule/sublevel_1/ren"
)

// DoPlay - 打つ
//
// * `command` - Example: `play black A19`
// ........................---- ----- ---
// ........................0    1     2
func (k *Kernel) DoPlay(command string, text_i_o i_text_i_o.ITextIO, log1 *logger.Logger) {
	var tokens = strings.Split(command, " ")
	var stoneName = tokens[1]

	var getDefaultStone = func() (bool, stone.Stone) {
		text_i_o.GoCommand(fmt.Sprintf("? unexpected stone:%s\n", stoneName))
		log1.J.Infow("error", "stone", stoneName)
		return false, stone.Stone_Space
	}

	var isOk1, stone = stone.GetStoneFromName(stoneName, getDefaultStone)
	if !isOk1 {
		return
	}

	var coord = tokens[2]
	// 着手点
	var placePlay = k.Position.Board.Coordinate.GetPointFromGtpMove(coord)

	// [O22o1o2o0] 石（または枠）の上に石を置こうとした
	var onMasonry = func() bool {
		text_i_o.GoCommand(fmt.Sprintf("? masonry my_stone:%s placePlay:%s\n", stone, k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error masonry", "my_stone", stone, "placePlay", k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o3o1o0] 相手の眼に石を置こうとした
	var onOpponentEye = func() bool {
		text_i_o.GoCommand(fmt.Sprintf("? opponent_eye my_stone:%s placePlay:%s\n", stone, k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error opponent_eye", "my_stone", stone, "placePlay", k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o4o1o0] 自分の眼に石を置こうとした
	var onForbiddenMyEye = func() bool {
		text_i_o.GoCommand(fmt.Sprintf("? my_eye my_stone:%s placePlay:%s\n", stone, k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error my_eye", "my_stone", stone, "placePlay", k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o7o2o0] コウに石を置こうとした
	var onKo = func() bool {
		text_i_o.GoCommand(fmt.Sprintf("? ko my_stone:%s placePlay:%s\n", stone, k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error ko", "my_stone", stone, "placePlay", k.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	var isOk = k.Play(stone, placePlay, log1,
		// [O22o1o2o0] ,onMasonry
		onMasonry,
		// [O22o3o1o0] ,onOpponentEye
		onOpponentEye,
		// [O22o4o1o0] ,onForbiddenMyEye
		onForbiddenMyEye,
		// [O22o7o2o0] ,onKo
		onKo)

	if isOk {
		text_i_o.GoCommand("=\n")
		log1.J.Infow("ok")
	}
}

// Play - 石を打つ
//
// Parameters
// ==========
// stoneA : Stone
// -
// placePlay : Point
// -
//
// Returns
// =======
// isOk : bool
// - 石を置けたら真、置けなかったら偽
func (k *Kernel) Play(stoneA stone.Stone, placePlay point.Point, logg *logger.Logger,
	// [O22o1o2o0] onMasonry
	onMasonry func() bool,
	// [O22o3o1o0] onOpponentEye
	onOpponentEye func() bool,
	// [O22o4o1o0] onForbiddenMyEye
	onForbiddenMyEye func() bool,
	// [O22o7o2o0] onKo
	onKo func() bool) bool {

	// [O22o1o2o0]
	if k.Position.Board.IsMasonry(placePlay) {
		return onMasonry()
	}

	// [O22o7o2o0] コウの判定
	if k.Record.IsKo(placePlay) {
		return onKo()
	}

	// [O22o6o1o0] Captured ルール
	var isExists4rensToRemove = false
	var o4rensToRemove [4]*rentype.Ren
	var isChecked4rensToRemove = false

	// [O22o3o1o0] 連と呼吸点の算出
	var renC, isFound = k.GetLiberty(placePlay)
	if isFound && renC.GetArea() == 1 { // 石Aを置いた交点を含む連Cについて、連Cの面積が1である（眼）
		if stoneA.GetColor() == renC.AdjacentColor.GetOpponent() {
			// かつ、連Cに隣接する連の色が、石Aのちょうど反対側の色であったなら、
			// 相手の眼に石を置こうとしたとみなす

			// [O22o6o1o0] 打ちあげる死に石の連を取得
			k.Position.Board.SetStoneAt(placePlay, stoneA) // いったん、石を置く
			isExists4rensToRemove, o4rensToRemove = k.GetRenToCapture(placePlay)
			isChecked4rensToRemove = true
			k.Position.Board.SetStoneAt(placePlay, stone.Stone_Space) // 石を取り除く

			if !isExists4rensToRemove {
				// `Captured` ルールと被らなければ
				return onOpponentEye()
			}

		} else if k.Position.CanNotPutOnMyEye && stoneA.GetColor() == renC.AdjacentColor {
			// [O22o4o1o0]
			// かつ、連Cに隣接する連の色が、石Aの色であったなら、
			// 自分の眼に石を置こうとしたとみなす
			return onForbiddenMyEye()

		}
	}

	// 石を置く
	k.Position.Board.SetStoneAt(placePlay, stoneA)

	// [O22o6o1o0] 打ちあげる死に石の連を取得
	if !isChecked4rensToRemove {
		isExists4rensToRemove, o4rensToRemove = k.GetRenToCapture(placePlay)
	}

	// [O22o7o2o0] コウの判定
	var capturedCount = 0 // アゲハマ

	// [O22o6o1o0] 死に石を打ちあげる
	if isExists4rensToRemove {
		for dir := 0; dir < 4; dir++ {
			var ren = o4rensToRemove[dir]

			if ren != nil {
				k.RemoveRen(ren)

				// [O22o7o2o0] コウの判定
				capturedCount += ren.GetArea()
			}
		}
	}

	// [O22o7o2o0] コウの判定
	var ko = point.Point(0)
	if capturedCount == 1 {
		ko = placePlay
	}

	// 棋譜に追加
	k.Record.Push(placePlay,
		// [O22o7o2o0] コウの判定
		ko)

	return true
}

// GetRenToCapture - 現在、着手後の盤面とする。打ち上げられる石の連を返却
//
// Returns
// -------
// isExists : bool
// renToRemove : [4]*Ren
// 隣接する東、北、西、南にある石を含む連
func (k *Kernel) GetRenToCapture(placePlay point.Point) (bool, [4]*rentype.Ren) {
	// [O22o6o1o0]
	var isExists bool
	var rensToRemove [4]*rentype.Ren
	var renIds = [4]point.Point{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	var setAdjacentPoint = func(dir board_coordinate.Cell_4Directions, adjacentP point.Point) {
		var adjacentR, isFound = k.GetLiberty(adjacentP)
		if isFound {
			// 同じ連を数え上げるのを防止する
			var renId = adjacentR.GetMinimumLocation()
			for i := board_coordinate.Cell_4Directions(0); i < dir; i++ {
				if renIds[i] == renId { // Idが既存
					return
				}
			}

			// 取れる石を見つけた
			if adjacentR.GetLibertyArea() < 1 {
				isExists = true
				rensToRemove[dir] = adjacentR
			}
		}
	}

	// 隣接する４方向
	k.Position.Board.ForeachNeumannNeighborhood(placePlay, setAdjacentPoint)

	return isExists, rensToRemove
}

// EOF [O19o0]
