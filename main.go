// BOF [O9o0]

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	dbg "github.com/muzudho/kifuwarabe-uec17/debugger"

	// Level 1
	logger "github.com/muzudho/kifuwarabe-uec17/kernel/level_1_for_maintenance/logger"

	// Level 2.1
	komi_float "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/komi_float"
	moves_num "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/moves_num"

	// Level 3.1
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/level_3_physical/sublevel_1/stone"

	// Level 4.1
	game_rule_settings "github.com/muzudho/kifuwarabe-uec17/kernel/level_4_game_rule/sublevel_1/game_rule_settings"

	// Level 31
	kernel_core "github.com/muzudho/kifuwarabe-uec17/kernel/level_31_controller"
)

// [O11o_1o0] グローバル変数として、バーチャルIOを１つ新規作成
// アプリケーションの中では 標準入出力は これを使うようにする
var virtualIo = dbg.NewVirtualIO()

func main() {
	// ========================================
	// flag - コマンドラインの解析
	// ========================================

	// 思考エンジン設定ファイル
	var (
		pEngineFilePath = flag.String("f", "engine.toml", "engine config file path")
		// [O11o__11o6o0] デバッグ用
		pIsDebug = flag.Bool("d", false, "for debug")
	)
	flag.Parse()
	// コマンドラインの第１引数で処理を振り分け
	var arg1 = flag.Arg(0)

	// ========================================
	// 思考エンジンの準備
	// ========================================

	// 思考エンジン設定ファイル
	var onError = func(err error) *Config {
		// ログファイルには出力できません。ログファイルはまだ読込んでいません

		// 強制終了
		panic(err)
	}
	var engineConfig = LoadEngineConfig(*pEngineFilePath, onError)

	// ========================================
	// 思考エンジンの準備　＞　ログ・ファイル
	// ========================================

	// ログファイル
	var plainTextLogFile, _ = os.OpenFile(engineConfig.GetPlainTextLog(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer plainTextLogFile.Close() // ログファイル使用済み時にファイルを閉じる
	// ログファイル
	var jsonLogFile, _ = os.OpenFile(engineConfig.GetJsonLog(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer jsonLogFile.Close() // ログファイル使用済み時にファイルを閉じる
	// カスタマイズしたロガーを使うなら
	var log1 = logger.NewSugaredLoggerForGame(plainTextLogFile, jsonLogFile)

	// [O11o__11o6o0] デバッグ用
	if *pIsDebug {
		virtualIo.ReplaceInputToFileLines("./debug.input.txt")
	}

	// ========================================
	// コマンドラインの第一引数で処理を振り分ける
	// ========================================

	switch arg1 {
	case "hello":
		// コマンドライン例： `go run . hello`
		fmt.Println("Hello, World!")

	case "welcome":
		// ロガーのテスト
		log1.C.Infof("Welcome! name:'%s' weight:%.1f x:%d", "nihon taro", 92.6, 3)
		log1.J.Infow("Welcome!",
			"name", "nihon taro", "weight", 92.6, "x", 3)

		// この上に分岐を挟んでいく
		// ---------------------

	default:
		// fmt.Println("go run . {programName}")

		// [O12o__11o_4o0] 棋譜の初期化に利用
		var onUnknownTurn = func() stone.Stone {
			var errMsg = fmt.Sprintf("? unexpected play_first:%s", engineConfig.GetPlayFirst())
			log1.C.Info(errMsg)
			log1.J.Infow("error", "play_first", engineConfig.GetPlayFirst())
			panic(errMsg)
		}

		// [O11o_3o0]
		var gameRuleSettings = game_rule_settings.NewGameRuleSettings(komi_float.KomiFloat(engineConfig.GetKomi()), moves_num.MovesNum(engineConfig.GetMaxPositionNumber()))
		var kernel1 = kernel_core.NewDirtyKernel(*gameRuleSettings, engineConfig.GetBoardSize(), engineConfig.GetBoardSize(),
			// [O12o__11o_4o0] 棋譜の初期化
			moves_num.MovesNum(engineConfig.GetMaxPositionNumber()),
			stone.GetStoneOrDefaultFromTurn(engineConfig.GetPlayFirst(), onUnknownTurn))
		// 設定ファイルの内容をカーネルへ反映
		kernel1.Position.Board.Init(engineConfig.GetBoardSize(), engineConfig.GetBoardSize())

		// [O11o_1o0] コンソール等からの文字列入力
		for virtualIo.ScannerScan() {
			var command = virtualIo.ScannerText()
			log1.C.Infof("# %s", command)             // 人間向けの出力
			log1.J.Infow("input", "command", command) // コンピューター向けの出力

			// [O11o_3o0]
			var isHandled = kernel1.Execute(command, log1)
			if isHandled {
				continue
			}

			// [O11o_1o0]
			var tokens = strings.Split(command, " ")
			switch tokens[0] {

			// この下にコマンドを挟んでいく
			// -------------------------

			case "quit": // [O11o_1o0]
				// os.Exit(0)
				return

			// この上にコマンドを挟んでいく
			// -------------------------

			default: // [O11o_1o0]
				log1.C.Infof("? unknown_command command:'%s'\n", tokens[0])
				log1.J.Infow("? unknown_command", "command", tokens[0])
			}
		}
	}
}

// EOF [O9o0]
