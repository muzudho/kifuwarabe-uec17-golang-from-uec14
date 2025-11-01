package level_31_controller

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	// Level 1
	logger "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_1_for_maintenance/logger"

	// Level 2.1
	geta "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_2_conceptual/sublevel_1/geta"
	moves_num "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_2_conceptual/sublevel_1/moves_num"
	point "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2.2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_2_conceptual/sublevel_2/board_coordinate"

	// Level 3.1
	stone "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_3_physical/sublevel_1/stone"

	// Level 4.1
	game_rule_settings "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_4_game_rule/sublevel_1/game_rule_settings"

	// Level 6.1
	record_item "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_6_database/sublevel_1/record_item"

	// Level 6.2
	record "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_6_database/sublevel_2/record"

	// Level 6.3
	ren_db "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_6_database/sublevel_3/ren_db"

	// Level 7.1
	position "github.com/muzudho/kifuwarabe-uec17-golang/kernel/level_7_misc/sublevel_1/position"
)

// Kernel - カーネル
type Kernel struct {
	// Position - 局面
	Position *position.Position

	// Record - [O12o__11o_3o0] 棋譜
	Record record.Record

	// RenDb - [O12o__11o__10o3o0] 連データベース
	renDb *ren_db.RenDb
}

// NewDirtyKernel - カーネルの新規作成
// - 一部のメンバーは、初期化されていないので、別途初期化処理が要る
func NewDirtyKernel(gameRuleSettings game_rule_settings.GameRuleSettings, boardWidht int, boardHeight int, maxMovesNum moves_num.MovesNum, playFirst stone.Stone) *Kernel {

	var k = new(Kernel)
	k.Position = position.NewDirtyPosition(gameRuleSettings, boardWidht, boardHeight)

	// [O12o__11o_2o0] 棋譜の初期化
	k.Record = *record.NewRecord(maxMovesNum, k.Position.Board.Coordinate.GetMemoryArea(), playFirst)

	// RenDb - [O12o__11o__10o3o0] 連データベース
	k.renDb = ren_db.NewRenDb(k.Position.Board.Coordinate.GetWidth(), k.Position.Board.Coordinate.GetHeight())

	return k
}

// ReadCommand - 実行
//
// Returns
// -------
// isHandled : bool
// 正常終了またはエラーなら真、無視したら偽
func (k *Kernel) ReadCommand(command string, log1 *logger.Logger) bool {

	var tokens = strings.Split(command, " ")
	switch tokens[0] {

	// ========================================
	// GTP 対応　＞　大会参加向け
	// ========================================

	// 盤サイズの設定
	// Example: `boardsize 19`
	case "boardsize": // [O15o__11o0]
		var sideLength, err = strconv.Atoi(tokens[1])

		if err != nil {
			log1.C.Infof("? unexpected sideLength:%s\n", tokens[1])
			log1.J.Infow("error", "sideLength", tokens[1])
			return true
		}

		k.Position.Board.Init(sideLength, sideLength)
		log1.C.Info("=\n")
		log1.J.Infow("ok")

		return true

	// 石を置く
	// Example: `play black A19`
	case "play":
		k.DoPlay(command, log1)
		return true

	// ========================================
	// 独自実装
	// ========================================

	// 独自実装：　盤をファイルから読み込んでセットする
	// Example: `board_set file data/board1.txt`
	case "board_set":
		k.DoSetBoard(command, log1)
		log1.C.Infof("=\n")
		log1.J.Infow("ok")
		return true

	// 独自実装：　人間向けに簡易の盤表示
	// Example: `board`
	case "board":
		{
			// 25列まで対応
			const fileSimbols = "ABCDEFGHJKLMNOPQRSTUVWXYZ"
			// 25行まで対応
			var rankSimbols = strings.Split("  , 1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25", ",")

			var filesMax = int(math.Min(25, float64(k.Position.Board.Coordinate.GetWidth())))
			var rowsMax = int(math.Min(25, float64(k.Position.Board.Coordinate.GetHeight())))
			var filesLabel = fileSimbols[:filesMax]

			var sb strings.Builder
			// 枠の上辺
			sb.WriteString(fmt.Sprintf(`= board:'''
.     %s
.    `, filesLabel))

			var rowNumber = 1
			var setPoint = func(point point.Point) {
				var stone = k.Position.Board.Cells[point]
				sb.WriteString(fmt.Sprintf("%v", stone))
			}
			var doNewline = func() {
				var rankLabel string
				if rowNumber <= rowsMax {
					rankLabel = rankSimbols[rowNumber]
				} else {
					rankLabel = ""
				}

				sb.WriteString(fmt.Sprintf("\n. %2s ", rankLabel))
				rowNumber++
			}
			k.Position.Board.GetCoordinate().ForeachLikeText(setPoint, doNewline)
			sb.WriteString("\n. '''\n")
			log1.C.Info(sb.String())
		}
		// コンピューター向けの出力
		{
			var sb strings.Builder

			var setPoint = func(point point.Point) {
				var stone = k.Position.Board.Cells[point]
				sb.WriteString(fmt.Sprintf("%v", stone))
			}
			var doNewline = func() {
				// pass
			}
			k.Position.Board.GetCoordinate().ForeachLikeText(setPoint, doNewline)
			log1.J.Infow("output", "board", sb.String())
		}
		return true

	// 独自実装：目に打たない
	// Example 1: `can_not_put_on_my_eye get`
	// Example 2: `can_not_put_on_my_eye set true``
	case "can_not_put_on_my_eye": // [O22o4o2o_1o0]
		var method = tokens[1]
		switch method {
		case "get":
			var value = k.Position.CanNotPutOnMyEye
			log1.C.Infof("= %t\n", value)
			log1.J.Infow("ok", "value", value)
			return true

		case "set":
			var value = tokens[2]
			switch value {
			case "true":
				k.Position.CanNotPutOnMyEye = true
				return true
			case "false":
				k.Position.CanNotPutOnMyEye = false
				return true
			default:
				log1.C.Infof("? unexpected method:%s value:%s\n", method, value)
				log1.J.Infow("error", "method", method, "value", value)
				return true
			}

		default:
			log1.C.Infof("? unexpected method:%s\n", method)
			log1.J.Infow("error", "method", method)
			return true
		}

	// 独自実装：　すべての連を取得
	// Example: `find_all_rens`
	case "find_all_rens":
		k.FindAllRens()
		log1.C.Infof("=\n")
		log1.J.Infow("ok")
		return true

	// Example: "record"
	case "record":
		var sb strings.Builder

		var setPoint = func(movesNum1 moves_num.MovesNum, item *record_item.RecordItem) {
			var positionNth = movesNum1 + geta.Geta // 基数を序数に変換
			var coord = k.Position.Board.Coordinate.GetGtpMoveFromPoint(item.PlacePlay)
			// sb.WriteString(fmt.Sprintf("[%d]%s ", positionNth, coord))

			// [O22o7o4o0] コウを追加
			var koStr string
			if item.Ko == point.Point(0) {
				koStr = ""
			} else {
				koStr = fmt.Sprintf("(%s)", k.Position.Board.Coordinate.GetGtpMoveFromPoint(item.Ko))
			}
			sb.WriteString(fmt.Sprintf("[%d]%s%s ", positionNth, coord, koStr))
		}

		k.Record.ForeachItem(setPoint)

		var text = sb.String()
		if 0 < len(text) {
			text = text[:len(text)-1]
		}
		log1.C.Infof("= record:'%s'\n", text)
		log1.J.Infow("ok", "record", text)
		return true

	// Example: `remove_ren B2`
	case "remove_ren":
		var coord = tokens[1]
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(coord)
		var ren, isFound = k.GetLiberty(point)
		if isFound {
			k.RemoveRen(ren)
			log1.C.Infof("=\n")
			log1.J.Infow("ok")
			return true
		}

		log1.C.Infof("? not found ren coord:%s%\n", coord)
		log1.J.Infow("error not found ren", "coord", coord)
		return false

	// 連データベースの内容をダンプ出力
	// Example: `rendb_dump`
	case "rendb_dump":
		var text = k.renDb.Dump()
		log1.C.Infof("= dump'''%s\n'''\n", text)
		log1.J.Infow("ok", "dump", text)
		return true

	// Example: `rendb_load data/ren_db1_temp.json`
	// * ファイルパスにスペースがはいっていてはいけない
	case "rendb_load":
		var path = tokens[1]
		var onError = func(err error) bool {
			log1.C.Infof("? error:%s\n", err)
			log1.J.Infow("error", "err", err)
			return false
		}

		var isOk = k.LoadRenDb(path, onError)
		if isOk {
			log1.C.Infof("=\n")
			log1.J.Infow("ok")
			return true
		}
		return false

	// Example: `rendb_save data/ren_db1_temp.json`
	// * ファイルパスにスペースがはいっていてはいけない
	case "rendb_save":
		var path = tokens[1]

		var convertLocation = func(location point.Point) string {
			return k.Position.Board.Coordinate.GetGtpMoveFromPoint(location)
		}

		var onError = func(err error) bool {
			log1.C.Infof("? error:%s\n", err)
			log1.J.Infow("error", "err", err)
			return false
		}

		var isOk = k.renDb.Save(path, convertLocation, onError)
		if isOk {
			log1.C.Infof("=\n")
			log1.J.Infow("ok")
			return true
		}

		return false

	// 交点の符号を整数に変換
	// Example: "test_coord A13"
	case "test_coord":
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
		log1.C.Infof("= %d\n", point)
		log1.J.Infow("output", "point", point)
		return true

	// FIXME: A を入れると A を返す？
	// Example: "test_file A"
	case "test_file":
		var file = board_coordinate.GetFileFromCode(tokens[1])
		log1.C.Infof("= %s\n", file)
		log1.J.Infow("output", "file", file)
		return true

	// Example: "test_get_liberty B2"
	case "test_get_liberty":
		var coord = tokens[1]
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(coord)
		var ren, isFound = k.GetLiberty(point)
		if isFound {
			log1.C.Infof("= ren stone:%s area:%d libertyArea:%d adjacentColor:%s\n", ren.Stone, ren.GetArea(), ren.GetLibertyArea(), ren.AdjacentColor)
			log1.J.Infow("output ren", "color", ren.Stone, "area", ren.GetArea(), "libertyArea", ren.GetLibertyArea(), "adjacentColor", ren.AdjacentColor)
			return true
		}

		log1.C.Infof("? not found ren coord:%s%\n", coord)
		log1.J.Infow("error not found ren", "coord", coord)
		return false

	// 交点の符号を整数に変換
	// Example: "test_get_point_from_code A1"
	case "test_get_point_from_code":
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
		var code = k.Position.Board.Coordinate.GetGtpMoveFromPoint(point)
		log1.C.Infof("= %d %s", point, code)
		log1.J.Infow("ok", "point", point, "code", code)
		return true

	// X, Y を 交点番号に変換
	// Example: "test_get_point_from_xy 2 3"
	case "test_get_point_from_xy":
		var x, errX = strconv.Atoi(tokens[1])
		if errX != nil {
			log1.C.Infof("? unexpected x:%s\n", tokens[1])
			log1.J.Infow("error", "x", tokens[1], "err", errX)
			return true
		}
		var y, errY = strconv.Atoi(tokens[2])
		if errY != nil {
			log1.C.Infof("? unexpected y:%s\n", tokens[2])
			log1.J.Infow("error", "y", tokens[2], "err", errY)
			return true
		}

		var point = k.Position.Board.Coordinate.GetPointFromXy(x, y)
		log1.C.Infof("= %d\n", point)
		log1.J.Infow("output", "point", point)
		return true

	// 段を返す？
	// Example: "test_rank 13"
	case "test_rank":
		var rank = board_coordinate.GetRankFromCode(tokens[1])
		log1.C.Infof("= %s\n", rank)
		log1.J.Infow("output", "rank", rank)
		return true

	// 整数を列符号に変換
	// Example: "test_x 18"
	case "test_x":
		var x, err = strconv.Atoi(tokens[1])
		if err != nil {
			log1.C.Infof("? unexpected x:%s\n", tokens[1])
			log1.J.Infow("error", "x", tokens[1])
			return true
		}
		var file = board_coordinate.GetFileFromX(x)
		log1.C.Infof("= %s\n", file)
		log1.J.Infow("output", "file", file)
		return true

	// 整数を段符号に変換
	// Example: "test_y 18"
	case "test_y":
		var y, err = strconv.Atoi(tokens[1])
		if err != nil {
			log1.C.Infof("? unexpected y:%s\n", tokens[1])
			log1.J.Infow("error", "y", tokens[1])
			return true
		}
		var rank = board_coordinate.GetRankFromY(y)
		log1.C.Infof("= %s\n", rank)
		log1.J.Infow("output", "rank", rank)
		return true

	default:
	}

	return false
}
