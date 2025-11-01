// BOF [O11o_3o0]

package kernel

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	// Level 1
	geta "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/geta"
	moves_num "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/moves_num"
	point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"

	// Level 2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/board_coordinate"
	game_rule_settings "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/game_rule_settings"
	record_item "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/record_item"
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/stone"

	// Level 3
	record "github.com/muzudho/kifuwarabe-uec17/kernel/types/level3/record"

	// Level 4
	position "github.com/muzudho/kifuwarabe-uec17/kernel/types/level4/position"
	ren_db "github.com/muzudho/kifuwarabe-uec17/kernel/types/level4/ren_db"
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

// Execute - 実行
//
// Returns
// -------
// isHandled : bool
// 正常終了またはエラーなら真、無視したら偽
func (k *Kernel) Execute(command string, logg *Logger) bool {

	var tokens = strings.Split(command, " ")
	switch tokens[0] {

	// この下にコマンドを挟んでいく
	// -------------------------

	case "board_set": // [O15o__14o2o0]
		// Example: `board_set file data/board1.txt`
		k.DoSetBoard(command, logg)
		logg.C.Infof("=\n")
		logg.J.Infow("ok")
		return true

	case "board": // [O13o0]
		// 人間向けの出力
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
			logg.C.Info(sb.String())
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
			logg.J.Infow("output", "board", sb.String())
		}
		return true

	case "boardsize": // [O15o__11o0]
		// Example: `boardsize 19`
		var sideLength, err = strconv.Atoi(tokens[1])

		if err != nil {
			logg.C.Infof("? unexpected sideLength:%s\n", tokens[1])
			logg.J.Infow("error", "sideLength", tokens[1])
			return true
		}

		k.Position.Board.Init(sideLength, sideLength)
		logg.C.Info("=\n")
		logg.J.Infow("ok")

		return true

	case "can_not_put_on_my_eye": // [O22o4o2o_1o0]
		// Example 1: "can_not_put_on_my_eye get"
		// Example 2: "can_not_put_on_my_eye set true"
		var method = tokens[1]
		switch method {
		case "get":
			var value = k.Position.CanNotPutOnMyEye
			logg.C.Infof("= %t\n", value)
			logg.J.Infow("ok", "value", value)
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
				logg.C.Infof("? unexpected method:%s value:%s\n", method, value)
				logg.J.Infow("error", "method", method, "value", value)
				return true
			}

		default:
			logg.C.Infof("? unexpected method:%s\n", method)
			logg.J.Infow("error", "method", method)
			return true
		}

	case "find_all_rens": // [O23o_2o2o0]
		// Example: `find_all_rens`
		k.FindAllRens()
		logg.C.Infof("=\n")
		logg.J.Infow("ok")
		return true

	case "play": // [O20o0]
		// Example: `play black A19`
		k.DoPlay(command, logg)
		return true

	case "record": // [O12o__11o_5o0]
		// Example: "record"
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
		logg.C.Infof("= record:'%s'\n", text)
		logg.J.Infow("ok", "record", text)
		return true

	case "remove_ren": // [O22o5o2o0]
		// Example: `remove_ren B2`
		var coord = tokens[1]
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(coord)
		var ren, isFound = k.GetLiberty(point)
		if isFound {
			k.RemoveRen(ren)
			logg.C.Infof("=\n")
			logg.J.Infow("ok")
			return true
		}

		logg.C.Infof("? not found ren coord:%s%\n", coord)
		logg.J.Infow("error not found ren", "coord", coord)
		return false

	case "rendb_dump": // [O12o__11o__10o4o0]
		var text = k.renDb.Dump()
		logg.C.Infof("= dump'''%s\n'''\n", text)
		logg.J.Infow("ok", "dump", text)
		return true

	case "rendb_load": // [O12o__11o__10o5o__10o1o0]
		// Example: `rendb_load data/ren_db1_temp.json`
		// * ファイルパスにスペースがはいっていてはいけない
		var path = tokens[1]
		var onError = func(err error) bool {
			logg.C.Infof("? error:%s\n", err)
			logg.J.Infow("error", "err", err)
			return false
		}

		var isOk = k.LoadRenDb(path, onError)
		if isOk {
			logg.C.Infof("=\n")
			logg.J.Infow("ok")
			return true
		}
		return false

	case "rendb_save": // [O12o__11o__10o4o0]
		// Example: `rendb_save data/ren_db1_temp.json`
		// * ファイルパスにスペースがはいっていてはいけない
		var path = tokens[1]

		var convertLocation = func(location point.Point) string {
			return k.Position.Board.Coordinate.GetGtpMoveFromPoint(location)
		}

		var onError = func(err error) bool {
			logg.C.Infof("? error:%s\n", err)
			logg.J.Infow("error", "err", err)
			return false
		}

		var isOk = k.renDb.Save(path, convertLocation, onError)
		if isOk {
			logg.C.Infof("=\n")
			logg.J.Infow("ok")
			return true
		}

		return false

	case "test_coord": // [O12o__10o2o0]
		// Example: "test_coord A13"
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
		logg.C.Infof("= %d\n", point)
		logg.J.Infow("output", "point", point)
		return true

	case "test_file": // [O12o__10o2o0]
		// Example: "test_file A"
		var file = board_coordinate.GetFileFromCode(tokens[1])
		logg.C.Infof("= %s\n", file)
		logg.J.Infow("output", "file", file)
		return true

	case "test_get_liberty": // [O22o2o5o0]
		// Example: "test_get_liberty B2"
		var coord = tokens[1]
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(coord)
		var ren, isFound = k.GetLiberty(point)
		if isFound {
			logg.C.Infof("= ren stone:%s area:%d libertyArea:%d adjacentColor:%s\n", ren.Stone, ren.GetArea(), ren.GetLibertyArea(), ren.AdjacentColor)
			logg.J.Infow("output ren", "color", ren.Stone, "area", ren.GetArea(), "libertyArea", ren.GetLibertyArea(), "adjacentColor", ren.AdjacentColor)
			return true
		}

		logg.C.Infof("? not found ren coord:%s%\n", coord)
		logg.J.Infow("error not found ren", "coord", coord)
		return false

	case "test_get_point_from_code": // [O16o1o0]
		// Example: "test_get_point_from_code A1"
		var point = k.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
		var code = k.Position.Board.Coordinate.GetGtpMoveFromPoint(point)
		logg.C.Infof("= %d %s", point, code)
		logg.J.Infow("ok", "point", point, "code", code)
		return true

	case "test_get_point_from_xy": // [O12o__11o2o0]
		// Example: "test_get_point_from_xy 2 3"
		var x, errX = strconv.Atoi(tokens[1])
		if errX != nil {
			logg.C.Infof("? unexpected x:%s\n", tokens[1])
			logg.J.Infow("error", "x", tokens[1], "err", errX)
			return true
		}
		var y, errY = strconv.Atoi(tokens[2])
		if errY != nil {
			logg.C.Infof("? unexpected y:%s\n", tokens[2])
			logg.J.Infow("error", "y", tokens[2], "err", errY)
			return true
		}

		var point = k.Position.Board.Coordinate.GetPointFromXy(x, y)
		logg.C.Infof("= %d\n", point)
		logg.J.Infow("output", "point", point)
		return true

	case "test_rank": // [O12o__10o2o0]
		// Example: "test_rank 13"
		var rank = board_coordinate.GetRankFromCode(tokens[1])
		logg.C.Infof("= %s\n", rank)
		logg.J.Infow("output", "rank", rank)
		return true

	case "test_x": // [O12o__10o2o0]
		// Example: "test_x 18"
		var x, err = strconv.Atoi(tokens[1])
		if err != nil {
			logg.C.Infof("? unexpected x:%s\n", tokens[1])
			logg.J.Infow("error", "x", tokens[1])
			return true
		}
		var file = board_coordinate.GetFileFromX(x)
		logg.C.Infof("= %s\n", file)
		logg.J.Infow("output", "file", file)
		return true

	case "test_y": // [O12o__10o2o0]
		// Example: "test_y 18"
		var y, err = strconv.Atoi(tokens[1])
		if err != nil {
			logg.C.Infof("? unexpected y:%s\n", tokens[1])
			logg.J.Infow("error", "y", tokens[1])
			return true
		}
		var rank = board_coordinate.GetRankFromY(y)
		logg.C.Infof("= %s\n", rank)
		logg.J.Infow("output", "rank", rank)
		return true

	// この上にコマンドを挟んでいく
	// -------------------------

	default:
	}

	return false
}

// EOF [O11o_3o0]
