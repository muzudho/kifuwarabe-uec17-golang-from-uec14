// BOF [O15o__14o1o0]

package core

import (
	"os"
	"strings"

	logger "github.com/muzudho/kifuwarabe-uec17/kernel/level_1_for_maintenance/logger"

	// Level 1
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/level_3_physical/sublevel_1/stone"
)

// DoSetBoard - 盤面を設定する
//
// コマンドラインの複数行入力は難しいので、ファイルから取ることにする
// * `command` - Example: `board_set file data/board1.txt`
// ........................--------- ---- ---------------
// ........................0         1    2
func (k *Kernel) DoSetBoard(command string, logg *logger.Logger) {
	var tokens = strings.Split(command, " ")

	if tokens[1] == "file" {
		var filePath = tokens[2]

		var fileData, err = os.ReadFile(filePath)
		if err != nil {
			logg.C.Infof("? unexpected file:%s\n", filePath)
			logg.J.Infow("error", "file", filePath)
			return
		}

		var getDefaultStone = func() (bool, stone.Stone) {
			return false, stone.Stone_Space
		}

		var size = k.Position.Board.Coordinate.GetMemoryArea()
		var i point.Point = 0
		for _, c := range string(fileData) {
			var str = string([]rune{c})
			var isOk, stone = stone.GetStoneFromChar(str, getDefaultStone)

			if isOk {
				if size <= int(i) {
					// 配列サイズ超過
					logg.C.Infof("? board out of bounds i:%d size:%d\n", i, size)
					logg.J.Infow("error board out of bounds", "i", i, "size", size)
					return
				}

				k.Position.Board.SetStoneAt(i, stone)
				i++
			}
		}

		// サイズが足りていないなら、エラー
		if int(i) != size {
			logg.C.Infof("? not enough size i:%d size:%d\n", i, size)
			logg.J.Infow("error not enough size", "i", i, "size", size)
			return
		}

		// [O23o_2o3o_1o0] 連データベース初期化
		k.renDb.Init(k.Position.Board.Coordinate.GetWidth(), k.Position.Board.Coordinate.GetHeight())
		k.FindAllRens()
	}
}
