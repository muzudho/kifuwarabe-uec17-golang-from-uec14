// BOF [O15o__14o1o0]

package kernel

import (
	"os"
	"strings"

	point "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/point"
	types2 "github.com/muzudho/kifuwarabe-uec17/kernel/types2"
)

// DoSetBoard - 盤面を設定する
//
// コマンドラインの複数行入力は難しいので、ファイルから取ることにする
// * `command` - Example: `board_set file data/board1.txt`
// ........................--------- ---- ---------------
// ........................0         1    2
func (k *Kernel) DoSetBoard(command string, logg *Logger) {
	var tokens = strings.Split(command, " ")

	if tokens[1] == "file" {
		var filePath = tokens[2]

		var fileData, err = os.ReadFile(filePath)
		if err != nil {
			logg.C.Infof("? unexpected file:%s\n", filePath)
			logg.J.Infow("error", "file", filePath)
			return
		}

		var getDefaultStone = func() (bool, types2.Stone) {
			return false, types2.Stone_Space
		}

		var size = k.Position.Board.coordinate.GetMemoryArea()
		var i point.Point = 0
		for _, c := range string(fileData) {
			var str = string([]rune{c})
			var isOk, stone = types2.GetStoneFromChar(str, getDefaultStone)

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
		k.renDb.Init(k.Position.Board.coordinate.GetWidth(), k.Position.Board.coordinate.GetHeight())
		k.FindAllRens()
	}
}
