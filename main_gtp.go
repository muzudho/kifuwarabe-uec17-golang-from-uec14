package main

import (
	"fmt"
	"math"
	"strconv"

	// Entities
	color "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/color"
	komi_float "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/komi_float"
	moves_num "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/moves_num"
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"

	//
	i_text_io "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/interfaces/part_1_facility/chapter_1_io/section_1/i_text_io"

	game_rule_settings "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_4_game_rule/sublevel_1/game_rule_settings"

	kernel_core "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_31_controller"

	logger "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_7_presenter/chapter_1_io/section_1/logger"

	"github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_1/geta"

	"strings"

	board_coordinate "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_2_conceptual/sublevel_2/board_coordinate"
	record_item "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_6_database/sublevel_1/record_item"
)

func LoopGTP(text_io1 i_text_io.ITextIO, log1 *logger.Logger, engineConfig *Config) {
	// [O12o__11o_4o0] 棋譜の初期化に利用
	var onUnknownTurn = func() color.Color {
		var errMsg = fmt.Sprintf("? unexpected play_first:%s", engineConfig.GetPlayFirst())
		text_io1.SendCommand(errMsg)
		log1.J.Infow("error", "play_first", engineConfig.GetPlayFirst())
		panic(errMsg)
	}

	// [O11o_3o0]
	var gameRuleSettings = game_rule_settings.NewGameRuleSettings(komi_float.KomiFloat(engineConfig.GetKomi()), moves_num.MovesNum(engineConfig.GetMaxPositionNumber()))
	var kernel1 = kernel_core.NewDirtyKernel(*gameRuleSettings, engineConfig.GetBoardSize(), engineConfig.GetBoardSize(),
		// [O12o__11o_4o0] 棋譜の初期化
		moves_num.MovesNum(engineConfig.GetMaxPositionNumber()),
		color.GetStoneOrDefaultFromTurn(engineConfig.GetPlayFirst(), onUnknownTurn))
	// 設定ファイルの内容をカーネルへ反映
	kernel1.Position.Board.Init(engineConfig.GetBoardSize(), engineConfig.GetBoardSize())

	// [O11o_1o0] コンソール等からの文字列入力
	for virtualIo.ScannerScan() {
		var command = virtualIo.ScannerText()
		text_io1.ReceivedCommand(command)

		var tokens = strings.Split(command, " ")
		switch tokens[0] {

		// ========================================
		// GTP 対応　＞　大会参加最低限
		// ========================================

		// 使用可能なコマンドのリスト
		// Example: `list_commands`
		case "list_commands":
			// 最初の１個は頭に "= " を付ける必要があってめんどくさいので先に出力
			text_io1.SendCommand("= list_commands\n")

			items := []string{
				// 終了コマンド
				"quit",
				// ハンドシェイク
				"protocol_version", "name", "version",
				// 対局設定
				"boardsize", "komi"}
			for _, item := range items {
				text_io1.SendCommand(fmt.Sprintf("%s\n", item))
			}

		// ========================================
		// GTP 対応　＞　大会参加最低限　＞　終了コマンド
		// ========================================

		case "quit": // [O11o_1o0]
			// os.Exit(0)
			return

		// ========================================
		// GTP 対応　＞　大会参加最低限　＞　ハンドシェイク
		// ========================================

		// 思考エンジンの名前
		// Example: `name`
		case "name":
			text_io1.SendCommand("= Kifuwarabe UEC17\n")

		// 思考エンジンのバージョン
		// Example: `version`
		case "version":
			text_io1.SendCommand("= 0.0.1\n")

		// プロトコルのバージョン
		// Example: `protocol_version`
		case "protocol_version":
			text_io1.SendCommand("= 2\n")

		// ========================================
		// GTP 対応　＞　大会参加最低限　＞　対局設定
		// ========================================

		// 盤サイズの設定
		// Example: `boardsize 19`
		case "boardsize":
			var sideLength, err = strconv.Atoi(tokens[1])

			if err != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected sideLength:%s\n", tokens[1]))
				log1.J.Infow("error", "sideLength", tokens[1])
				continue
			}

			kernel1.Position.Board.Init(sideLength, sideLength)
			text_io1.SendCommand("=\n")
			log1.J.Infow("ok")

		// コミの設定
		// Example: `komi 6.5`
		case "komi":
			komi, err := strconv.ParseFloat(tokens[1], 64)

			if err != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected komi:%s\n", tokens[1]))
				log1.J.Infow("error", "komi", tokens[1])
			}

			kernel1.Position.Board.GameRuleSettings.Komi = komi_float.KomiFloat(komi)
			text_io1.SendCommand("=\n")
			log1.J.Infow("ok")

		// ========================================
		// GTP 対応　＞　大会参加最低限　＞　対局
		// ========================================

		case "clear_board":
			kernel1.Position.Board.Init(engineConfig.GetBoardSize(), engineConfig.GetBoardSize())
			text_io1.SendCommand("= \n\n")

		// 石を置く
		// Example: `play black A19`
		case "play":
			kernel1.DoPlay(command, text_io1, log1)

		case "undo":
			// 未実装
			text_io1.SendCommand("= \n\n")

		case "genmove":
			// genmove black
			// genmove white
			// var color e.Stone
			// if 1 < len(tokens) && strings.ToLower(tokens[1][0:1]) == "w" {
			// 	color = 2
			// } else {
			// 	color = 1
			// }
			// var z = PlayComputerMoveLesson09a(position, color)
			// text_io1.SendCommand(fmt.Sprintf("= %s\n\n", p.GetGtpZ(position, z)))

		// ========================================
		// 独自コマンド
		// ========================================

		// 独自実装：　盤をファイルから読み込んでセットする
		// Example: `board_set file data/board1.txt`
		case "board_set":
			kernel1.DoSetBoard(command, text_io1, log1)
			text_io1.SendCommand("=\n")
			log1.J.Infow("ok")

		// 独自実装：　人間向けに簡易の盤表示
		// Example: `board`
		case "board":
			{
				// 25列まで対応
				const fileSimbols = "ABCDEFGHJKLMNOPQRSTUVWXYZ"
				// 25行まで対応
				var rankSimbols = strings.Split("  , 1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25", ",")

				var filesMax = int(math.Min(25, float64(kernel1.Position.Board.Coordinate.GetWidth())))
				var rowsMax = int(math.Min(25, float64(kernel1.Position.Board.Coordinate.GetHeight())))
				var filesLabel = fileSimbols[:filesMax]

				var sb strings.Builder
				// 枠の上辺
				sb.WriteString(fmt.Sprintf(`= board:'''
.     %s
.    `, filesLabel))

				var rowNumber = 1
				var setPoint = func(point point.Point) {
					var stone = kernel1.Position.Board.Cells[point]
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
				kernel1.Position.Board.GetCoordinate().ForeachLikeText(setPoint, doNewline)
				sb.WriteString("\n. '''\n")
				text_io1.SendCommand(sb.String())
			}
			// コンピューター向けの出力
			{
				var sb strings.Builder

				var setPoint = func(point point.Point) {
					var stone = kernel1.Position.Board.Cells[point]
					sb.WriteString(fmt.Sprintf("%v", stone))
				}
				var doNewline = func() {
					// pass
				}
				kernel1.Position.Board.GetCoordinate().ForeachLikeText(setPoint, doNewline)
				log1.J.Infow("output", "board", sb.String())
			}

		// 独自実装：目に打たない
		// Example 1: `can_not_put_on_my_eye get`
		// Example 2: `can_not_put_on_my_eye set true``
		case "can_not_put_on_my_eye": // [O22o4o2o_1o0]
			var method = tokens[1]
			switch method {
			case "get":
				var value = kernel1.Position.CanNotPutOnMyEye
				text_io1.SendCommand(fmt.Sprintf("= %t\n", value))
				log1.J.Infow("ok", "value", value)

			case "set":
				var value = tokens[2]
				switch value {
				case "true":
					kernel1.Position.CanNotPutOnMyEye = true
				case "false":
					kernel1.Position.CanNotPutOnMyEye = false
				default:
					text_io1.SendCommand(fmt.Sprintf("? unexpected method:%s value:%s\n", method, value))
					log1.J.Infow("error", "method", method, "value", value)
				}

			default:
				text_io1.SendCommand(fmt.Sprintf("? unexpected method:%s\n", method))
				log1.J.Infow("error", "method", method)
			}

		// 独自実装：　すべての連を取得
		// Example: `find_all_rens`
		case "find_all_rens":
			kernel1.FindAllRens()
			text_io1.SendCommand("=\n")
			log1.J.Infow("ok")

		// Example: "record"
		case "record":
			var sb strings.Builder

			var setPoint = func(movesNum1 moves_num.MovesNum, item *record_item.RecordItem) {
				var positionNth = movesNum1 + geta.Geta // 基数を序数に変換
				var coord = kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(item.PlacePlay)
				// sb.WriteString(fmt.Sprintf("[%d]%s ", positionNth, coord))

				// [O22o7o4o0] コウを追加
				var koStr string
				if item.Ko == point.Point(0) {
					koStr = ""
				} else {
					koStr = fmt.Sprintf("(%s)", kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(item.Ko))
				}
				sb.WriteString(fmt.Sprintf("[%d]%s%s ", positionNth, coord, koStr))
			}

			kernel1.Record.ForeachItem(setPoint)

			var text = sb.String()
			if 0 < len(text) {
				text = text[:len(text)-1]
			}
			text_io1.SendCommand(fmt.Sprintf("= record:'%s'\n", text))
			log1.J.Infow("ok", "record", text)

		// Example: `remove_ren B2`
		case "remove_ren":
			var coord = tokens[1]
			var point = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(coord)
			var ren, isFound = kernel1.GetLiberty(point)
			if isFound {
				kernel1.RemoveRen(ren)
				text_io1.SendCommand("=\n")
				log1.J.Infow("ok")
				continue
			}

			text_io1.SendCommand(fmt.Sprintf("? not found ren coord:%s\n", coord))
			log1.J.Infow("error not found ren", "coord", coord)

		// 連データベースの内容をダンプ出力
		// Example: `rendb_dump`
		case "rendb_dump":
			var text = kernel1.RenDb.Dump()
			text_io1.SendCommand(fmt.Sprintf("= dump'''%s\n'''\n", text))
			log1.J.Infow("ok", "dump", text)

		// Example: `rendb_load data/ren_db1_temp.json`
		// * ファイルパスにスペースがはいっていてはいけない
		case "rendb_load":
			var path = tokens[1]
			var onError = func(err error) bool {
				text_io1.SendCommand(fmt.Sprintf("? error:%s\n", err))
				log1.J.Infow("error", "err", err)
				return false
			}

			var isOk = kernel1.LoadRenDb(path, onError)
			if isOk {
				text_io1.SendCommand("=\n")
				log1.J.Infow("ok")
			}

		// Example: `rendb_save data/ren_db1_temp.json`
		// * ファイルパスにスペースがはいっていてはいけない
		case "rendb_save":
			var path = tokens[1]

			var convertLocation = func(location point.Point) string {
				return kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(location)
			}

			var onError = func(err error) bool {
				text_io1.SendCommand(fmt.Sprintf("? error:%s\n", err))
				log1.J.Infow("error", "err", err)
				return false
			}

			var isOk = kernel1.RenDb.Save(path, convertLocation, onError)
			if isOk {
				text_io1.SendCommand("=\n")
				log1.J.Infow("ok")
			}

		// 交点の符号を整数に変換
		// Example: "test_coord A13"
		case "test_coord":
			var point = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
			text_io1.SendCommand(fmt.Sprintf("= %d\n", point))
			log1.J.Infow("output", "point", point)

		// FIXME: A を入れると A を返す？
		// Example: "test_file A"
		case "test_file":
			var file = board_coordinate.GetFileFromCode(tokens[1])
			text_io1.SendCommand(fmt.Sprintf("= %s\n", file))
			log1.J.Infow("output", "file", file)

		// Example: "test_get_liberty B2"
		case "test_get_liberty":
			var coord = tokens[1]
			var point = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(coord)
			var ren, isFound = kernel1.GetLiberty(point)
			if isFound {
				text_io1.SendCommand(fmt.Sprintf("= ren stone:%s area:%d libertyArea:%d adjacentColor:%s\n", ren.Stone, ren.GetArea(), ren.GetLibertyArea(), ren.AdjacentColor))
				log1.J.Infow("output ren", "color", ren.Stone, "area", ren.GetArea(), "libertyArea", ren.GetLibertyArea(), "adjacentColor", ren.AdjacentColor)
				continue
			}

			text_io1.SendCommand(fmt.Sprintf("? not found ren coord:%s\n", coord))
			log1.J.Infow("error not found ren", "coord", coord)

		// 交点の符号を整数に変換
		// Example: "test_get_point_from_code A1"
		case "test_get_point_from_code":
			var point = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(tokens[1])
			var code = kernel1.Position.Board.Coordinate.GetGtpMoveFromPoint(point)
			text_io1.SendCommand(fmt.Sprintf("= %d %s", point, code))
			log1.J.Infow("ok", "point", point, "code", code)

		// X, Y を 交点番号に変換
		// Example: "test_get_point_from_xy 2 3"
		case "test_get_point_from_xy":
			var x, errX = strconv.Atoi(tokens[1])
			if errX != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected x:%s\n", tokens[1]))
				log1.J.Infow("error", "x", tokens[1], "err", errX)
				continue
			}

			var y, errY = strconv.Atoi(tokens[2])
			if errY != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected y:%s\n", tokens[2]))
				log1.J.Infow("error", "y", tokens[2], "err", errY)
				continue
			}

			var point = kernel1.Position.Board.Coordinate.GetPointFromXy(x, y)
			text_io1.SendCommand(fmt.Sprintf("= %d\n", point))
			log1.J.Infow("output", "point", point)

		// 段を返す？
		// Example: "test_rank 13"
		case "test_rank":
			var rank = board_coordinate.GetRankFromCode(tokens[1])
			text_io1.SendCommand(fmt.Sprintf("= %s\n", rank))
			log1.J.Infow("output", "rank", rank)

		// 整数を列符号に変換
		// Example: "test_x 18"
		case "test_x":
			var x, err = strconv.Atoi(tokens[1])
			if err != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected x:%s\n", tokens[1]))
				log1.J.Infow("error", "x", tokens[1])
				continue
			}

			var file = board_coordinate.GetFileFromX(x)
			text_io1.SendCommand(fmt.Sprintf("= %s\n", file))
			log1.J.Infow("output", "file", file)

		// 整数を段符号に変換
		// Example: "test_y 18"
		case "test_y":
			var y, err = strconv.Atoi(tokens[1])
			if err != nil {
				text_io1.SendCommand(fmt.Sprintf("? unexpected y:%s\n", tokens[1]))
				log1.J.Infow("error", "y", tokens[1])
				continue
			}

			var rank = board_coordinate.GetRankFromY(y)
			text_io1.SendCommand(fmt.Sprintf("= %s\n", rank))
			log1.J.Infow("output", "rank", rank)

		// ========================================
		// 未定義のコマンド
		// ========================================

		default:
			text_io1.SendCommand(fmt.Sprintf("? unknown_command command:'%s'\n", tokens[0]))
			log1.J.Infow("? unknown_command", "command", tokens[0])
		}
	}
}
