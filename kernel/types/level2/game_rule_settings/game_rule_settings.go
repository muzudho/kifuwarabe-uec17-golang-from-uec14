package kernel

import (
	// Level 1
	moves_num "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/moves_num"
	komi_float "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/komi_float"
)

// GameRuleSettings - 対局ルール設定
type GameRuleSettings struct {
	// コミ。 6.5 といった数字を入れるだけ。実行速度優先で 64bitに
	komi komi_float.KomiFloat

	// 上限手数
	maxMovesNum moves_num.MovesNum
}

// NewGameRuleSettings - 新規作成
func NewGameRuleSettings(komi komi_float.KomiFloat, maxMovesNum moves_num.MovesNum) *GameRuleSettings {
	var gr = new(GameRuleSettings)

	gr.komi = komi
	gr.maxMovesNum = maxMovesNum

	return gr
}

// GetKomi - コミ取得
func (gr *GameRuleSettings) GetKomi() komi_float.KomiFloat {
	return gr.komi
}

// GetMaxPositionNumber - 上限手数
func (gr *GameRuleSettings) GetMaxPositionNumber() moves_num.MovesNum {
	return gr.maxMovesNum
}
