package kernel

// Level 1
import (
	komi_float "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/komi_float"
	moves_num "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/moves_num"
)

// GameRule - 対局ルール
type GameRule struct {
	// コミ。 6.5 といった数字を入れるだけ。実行速度優先で 64bitに
	komi komi_float.KomiFloat

	// 上限手数
	maxMovesNum moves_num.MovesNum
}

// NewGameRule - 新規作成
func NewGameRule(komi komi_float.KomiFloat, maxMovesNum moves_num.MovesNum) *GameRule {
	var gr = new(GameRule)

	gr.komi = komi
	gr.maxMovesNum = maxMovesNum

	return gr
}

// GetKomi - コミ取得
func (gr *GameRule) GetKomi() komi_float.KomiFloat {
	return gr.komi
}

// GetMaxPositionNumber - 上限手数
func (gr *GameRule) GetMaxPositionNumber() moves_num.MovesNum {
	return gr.maxMovesNum
}
