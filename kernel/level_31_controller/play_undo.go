// BOF [O23o1o0]

package core

import logger "github.com/muzudho/kifuwarabe-uec17/kernel/level_1_for_maintenance/logger"

// DoUndoPlay - 石を打ったのを戻す
//
// * `command` - Example: `undo`
// ........................----
// ........................0
func (k *Kernel) DoUndoPlay(command string, logg *logger.Logger) {
	k.UndoPlay()
}

// UndoPlay - 石を打ったのを戻す
//
// Returns
// -------
// isOk : bool
//
//	石を置けたら真、置けなかったら偽
func (k *Kernel) UndoPlay() bool {

	// 初期局面から前には戻せない
	if k.Record.GetMovesNum() < 1 {
		return false
	}

	// TODO 置いた石を盤上から消す
	// TODO アゲハマを取る前の盤上の場所を棋譜に記録しておく。連単位？
	// TODO アゲハマを盤上に戻す
	// TODO 棋譜から一手消す。カレントも減らす

	return false
}

// EOF [O23o1o0]
