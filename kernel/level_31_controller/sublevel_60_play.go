// BOF [O19o0]

package level_31_controller

import (
	"fmt"
	"math"
	"strings"

	// Entities
	color "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/color"
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"

	// Section 1.1.1
	logger "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_7_presenter/chapter_1_io/section_1/logger"
	i_text_io "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/interfaces/part_1_facility/chapter_1_io/section_1/i_text_io"

	// Level 2.2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_2/board_coordinate"

	// Level 4.1
	rentype "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_4_game_rule/sublevel_1/ren"
)

// DoPlay - 打つ
//
// * `command` - Example: `play black A19`
// ........................---- ----- ---
// ........................0    1     2
func (kernel1 *Kernel) DoPlay(command string, text_io i_text_io.ITextIO, log1 *logger.Logger) {
	var tokens = strings.Split(command, " ")
	var stoneName = tokens[1]

	var getDefaultColor = func() (bool, color.Color) {
		text_io.SendCommand(fmt.Sprintf("? unexpected stone:%s\n", stoneName))
		log1.J.Infow("error", "stone", stoneName)
		return false, color.None
	}

	var isOk1, stone = color.GetColorFromCode(stoneName, getDefaultColor)
	if !isOk1 {
		return
	}

	var coord = tokens[2]
	// 着手点
	var placePlay = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(coord)

	// [O22o1o2o0] 石（または枠）の上に石を置こうとした
	var onMasonry = func() bool {
		text_io.SendCommand(fmt.Sprintf("? masonry my_stone:%s placePlay:%s\n", stone, kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error masonry", "my_stone", stone, "placePlay", kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o3o1o0] 相手の眼に石を置こうとした
	var onOpponentEye = func() bool {
		text_io.SendCommand(fmt.Sprintf("? opponent_eye my_stone:%s placePlay:%s\n", stone, kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error opponent_eye", "my_stone", stone, "placePlay", kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o4o1o0] 自分の眼に石を置こうとした
	var onForbiddenMyEye = func() bool {
		text_io.SendCommand(fmt.Sprintf("? my_eye my_stone:%s placePlay:%s\n", stone, kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error my_eye", "my_stone", stone, "placePlay", kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o7o2o0] コウに石を置こうとした
	var onKo = func() bool {
		text_io.SendCommand(fmt.Sprintf("? ko my_stone:%s placePlay:%s\n", stone, kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay)))
		log1.J.Infow("error ko", "my_stone", stone, "placePlay", kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	var isOk = kernel1.Play(stone, placePlay, log1,
		// [O22o1o2o0] ,onMasonry
		onMasonry,
		// [O22o3o1o0] ,onOpponentEye
		onOpponentEye,
		// [O22o4o1o0] ,onForbiddenMyEye
		onForbiddenMyEye,
		// [O22o7o2o0] ,onKo
		onKo)

	if isOk {
		text_io.SendCommand("=\n")
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
func (kernel1 *Kernel) Play(color1 color.Color, placePlay point.Point, logg *logger.Logger,
	// [O22o1o2o0] onMasonry
	onMasonry func() bool,
	// [O22o3o1o0] onOpponentEye
	onOpponentEye func() bool,
	// [O22o4o1o0] onForbiddenMyEye
	onForbiddenMyEye func() bool,
	// [O22o7o2o0] onKo
	onKo func() bool) bool {

	// [O22o1o2o0]
	if kernel1.Position.Board.IsMasonry(placePlay) {
		return onMasonry()
	}

	// [O22o7o2o0] コウの判定
	if kernel1.Record.IsKo(placePlay) {
		return onKo()
	}

	// [O22o6o1o0] Captured ルール
	var isExists4rensToRemove = false
	var o4rensToRemove [4]*rentype.Ren
	var isChecked4rensToRemove = false

	// [O22o3o1o0] 連と呼吸点の算出
	var renC, isFound = kernel1.GetLiberty(placePlay)
	if isFound && renC.GetArea() == 1 { // 石Aを置いた交点を含む連Cについて、連Cの面積が1である（眼）
		if color1 == renC.AdjacentColor.GetOpponent() {
			// かつ、連Cに隣接する連の色が、石Aのちょうど反対側の色であったなら、
			// 相手の眼に石を置こうとしたとみなす

			// [O22o6o1o0] 打ちあげる死に石の連を取得
			kernel1.Position.Board.SetStoneAt(placePlay, color1) // いったん、石を置く
			isExists4rensToRemove, o4rensToRemove = kernel1.GetRenToCapture(placePlay)
			isChecked4rensToRemove = true
			kernel1.Position.Board.SetStoneAt(placePlay, color.None) // 石を取り除く

			if !isExists4rensToRemove {
				// `Captured` ルールと被らなければ
				return onOpponentEye()
			}

		} else if kernel1.Position.CanNotPutOnMyEye && color1 == renC.AdjacentColor {
			// [O22o4o1o0]
			// かつ、連Cに隣接する連の色が、石Aの色であったなら、
			// 自分の眼に石を置こうとしたとみなす
			return onForbiddenMyEye()

		}
	}

	// 石を置く
	kernel1.Position.Board.SetStoneAt(placePlay, color1)

	// [O22o6o1o0] 打ちあげる死に石の連を取得
	if !isChecked4rensToRemove {
		isExists4rensToRemove, o4rensToRemove = kernel1.GetRenToCapture(placePlay)
	}

	// [O22o7o2o0] コウの判定
	var capturedCount = 0 // アゲハマ

	// [O22o6o1o0] 死に石を打ちあげる
	if isExists4rensToRemove {
		for dir := 0; dir < 4; dir++ {
			var ren = o4rensToRemove[dir]

			if ren != nil {
				kernel1.RemoveRen(ren)

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
	kernel1.Record.Push(placePlay,
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
func (kernel1 *Kernel) GetRenToCapture(placePlay point.Point) (bool, [4]*rentype.Ren) {
	// [O22o6o1o0]
	var isExists bool
	var rensToRemove [4]*rentype.Ren
	var renIds = [4]point.Point{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	var setAdjacentPoint = func(dir board_coordinate.Cell_4Directions, adjacentP point.Point) {
		var adjacentR, isFound = kernel1.GetLiberty(adjacentP)
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
	kernel1.Position.Board.ForeachNeumannNeighborhood(placePlay, setAdjacentPoint)

	return isExists, rensToRemove
}

// EOF [O19o0]
