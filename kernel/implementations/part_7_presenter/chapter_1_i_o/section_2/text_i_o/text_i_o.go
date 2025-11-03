package text_io

// Section 1.1.1
import (
	"fmt"

	logger "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_7_presenter/chapter_1_i_o/section_1/logger"
)

// TextIO - テキスト入出力
type TextIO struct {
	// ロガー
	log1 *logger.Logger
}

func NewTextIO(log1 *logger.Logger) *TextIO {
	var t = new(TextIO)
	t.log1 = log1
	return t
}

func (t *TextIO) GoCommand(command string) {
	fmt.Print(command)
	//t.log1.C.Info(command)
}
