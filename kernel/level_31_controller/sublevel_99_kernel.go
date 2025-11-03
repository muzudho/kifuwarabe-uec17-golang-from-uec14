package level_31_controller

import (
	// Entities
	moves_num "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/moves_num"
	stone "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/stone"

	// Level 4.1
	game_rule_settings "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_4_game_rule/sublevel_1/game_rule_settings"

	// Level 6.2
	record "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_6_database/sublevel_2/record"

	// Level 6.3
	ren_db "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_6_database/sublevel_3/ren_db"

	// Level 7.1
	position "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_7_misc/sublevel_1/position"
)

// Kernel - カーネル
type Kernel struct {
	// Position - 局面
	Position *position.Position

	// Record - [O12o__11o_3o0] 棋譜
	Record record.Record

	// RenDb - [O12o__11o__10o3o0] 連データベース
	RenDb *ren_db.RenDb
}

// NewDirtyKernel - カーネルの新規作成
// - 一部のメンバーは、初期化されていないので、別途初期化処理が要る
func NewDirtyKernel(gameRuleSettings game_rule_settings.GameRuleSettings, boardWidht int, boardHeight int, maxMovesNum moves_num.MovesNum, playFirst stone.Stone) *Kernel {

	var k = new(Kernel)
	k.Position = position.NewDirtyPosition(gameRuleSettings, boardWidht, boardHeight)

	// [O12o__11o_2o0] 棋譜の初期化
	k.Record = *record.NewRecord(maxMovesNum, k.Position.Board.Coordinate.GetMemoryArea(), playFirst)

	// RenDb - [O12o__11o__10o3o0] 連データベース
	k.RenDb = ren_db.NewRenDb(k.Position.Board.Coordinate.GetWidth(), k.Position.Board.Coordinate.GetHeight())

	return k
}
