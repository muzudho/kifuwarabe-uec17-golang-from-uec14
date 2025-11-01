// BOF [O19o0]

package kernel

import (
	"math"
	"strings"

	types1 "github.com/muzudho/kifuwarabe-uec17/kernel/types1"
	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
	types3 "github.com/muzudho/kifuwarabe-uec17/kernel/types3"
)

// DoPlay - 打つ
//
// * `command` - Example: `play black A19`
// ........................---- ----- ---
// ........................0    1     2
func (k *Kernel) DoPlay(command string, logg *Logger) {
	var tokens = strings.Split(command, " ")
	var stoneName = tokens[1]

	var getDefaultStone = func() (bool, types2.Stone) {
		logg.C.Infof("? unexpected stone:%s\n", stoneName)
		logg.J.Infow("error", "stone", stoneName)
		return false, types2.Stone_Space
	}

	var isOk1, stone = types2.GetStoneFromName(stoneName, getDefaultStone)
	if !isOk1 {
		return
	}

	var coord = tokens[2]
	// 着手点
	var placePlay = k.Position.Board.coordinate.GetPointFromGtpMove(coord)

	// [O22o1o2o0] 石（または枠）の上に石を置こうとした
	var onMasonry = func() bool {
		logg.C.Infof("? masonry my_stone:%s placePlay:%s\n", stone, k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		logg.J.Infow("error masonry", "my_stone", stone, "placePlay", k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o3o1o0] 相手の眼に石を置こうとした
	var onOpponentEye = func() bool {
		logg.C.Infof("? opponent_eye my_stone:%s placePlay:%s\n", stone, k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		logg.J.Infow("error opponent_eye", "my_stone", stone, "placePlay", k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o4o1o0] 自分の眼に石を置こうとした
	var onForbiddenMyEye = func() bool {
		logg.C.Infof("? my_eye my_stone:%s placePlay:%s\n", stone, k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		logg.J.Infow("error my_eye", "my_stone", stone, "placePlay", k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	// [O22o7o2o0] コウに石を置こうとした
	var onKo = func() bool {
		logg.C.Infof("? ko my_stone:%s placePlay:%s\n", stone, k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		logg.J.Infow("error ko", "my_stone", stone, "placePlay", k.Position.Board.coordinate.GetGtpMoveFromPoint(placePlay))
		return false
	}

	var isOk = k.Play(stone, placePlay, logg,
		// [O22o1o2o0] ,onMasonry
		onMasonry,
		// [O22o3o1o0] ,onOpponentEye
		onOpponentEye,
		// [O22o4o1o0] ,onForbiddenMyEye
		onForbiddenMyEye,
		// [O22o7o2o0] ,onKo
		onKo)

	if isOk {
		logg.C.Info("=\n")
		logg.J.Infow("ok")
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
func (k *Kernel) Play(stoneA types2.Stone, placePlay types1.Point, logg *Logger,
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
	var o4rensToRemove [4]*types3.Ren
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
			k.Position.Board.SetStoneAt(placePlay, types2.Stone_Space) // 石を取り除く

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
	var ko = types1.Point(0)
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
func (k *Kernel) GetRenToCapture(placePlay types1.Point) (bool, [4]*types3.Ren) {
	// [O22o6o1o0]
	var isExists bool
	var rensToRemove [4]*types3.Ren
	var renIds = [4]types1.Point{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	var setAdjacentPoint = func(dir types2.Cell_4Directions, adjacentP types1.Point) {
		var adjacentR, isFound = k.GetLiberty(adjacentP)
		if isFound {
			// 同じ連を数え上げるのを防止する
			var renId = adjacentR.GetMinimumLocation()
			for i := types2.Cell_4Directions(0); i < dir; i++ {
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
