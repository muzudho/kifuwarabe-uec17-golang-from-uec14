// BOF [O15o__14o1o0]

package level_31_controller

import (
	"fmt"
	"os"
	"strings"

	// Entities
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"
	stone "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/stone"

	// Section 1.1.1
	logger "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_7_presenter/chapter_1_io/section_1/logger"
	i_text_io "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/interfaces/part_1_facility/chapter_1_io/section_1/i_text_io"
)

// DoSetBoard - 盤面を設定する
//
// コマンドラインの複数行入力は難しいので、ファイルから取ることにする
// * `command` - Example: `board_set file data/board1.txt`
// ........................--------- ---- ---------------
// ........................0         1    2
func (kernel1 *Kernel) DoSetBoard(command string, text_io i_text_io.ITextIO, logg *logger.Logger) {
	var tokens = strings.Split(command, " ")

	if tokens[1] == "file" {
		var filePath = tokens[2]

		var fileData, err = os.ReadFile(filePath)
		if err != nil {
			text_io.SendCommand(fmt.Sprintf("? unexpected file:%s\n", filePath))
			logg.J.Infow("error", "file", filePath)
			return
		}

		var getDefaultStone = func() (bool, stone.Stone) {
			return false, stone.None
		}

		var size = kernel1.Position.Board.Coordinate.GetMemoryArea()
		var i point.Point = 0
		for _, c := range string(fileData) {
			var str = string([]rune{c})
			var isOk, stone = stone.GetStoneFromChar(str, getDefaultStone)

			if isOk {
				if size <= int(i) {
					// 配列サイズ超過
					text_io.SendCommand(fmt.Sprintf("? board out of bounds i:%d size:%d\n", i, size))
					logg.J.Infow("error board out of bounds", "i", i, "size", size)
					return
				}

				kernel1.Position.Board.SetStoneAt(i, stone)
				i++
			}
		}

		// サイズが足りていないなら、エラー
		if int(i) != size {
			text_io.SendCommand(fmt.Sprintf("? not enough size i:%d size:%d\n", i, size))
			logg.J.Infow("error not enough size", "i", i, "size", size)
			return
		}

		// [O23o_2o3o_1o0] 連データベース初期化
		kernel1.RenDb.Init(kernel1.Position.Board.Coordinate.GetWidth(), kernel1.Position.Board.Coordinate.GetHeight())
		kernel1.FindAllRens()
	}
}
