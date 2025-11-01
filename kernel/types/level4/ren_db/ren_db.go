package ren_db

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// Level 1
	moves_num "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/moves_num"
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"
	geta "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/geta"
	ren_id "github.com/muzudho/kifuwarabe-uec17/kernel/types/level1/ren_id"

	// Level 2
	board_coordinate "github.com/muzudho/kifuwarabe-uec17/kernel/types/level2/board_coordinate"

	// Level 3
	rentype "github.com/muzudho/kifuwarabe-uec17/kernel/types/level3/ren"
	ren_db_doc_header "github.com/muzudho/kifuwarabe-uec17/kernel/types/level3/ren_db_doc_header"
)

// RenDb - 連データベース
type RenDb struct {
	// Header - ヘッダー
	Header ren_db_doc_header.RenDbDocHeader `json:"header"`

	// 要素
	Rens map[ren_id.RenId]*rentype.Ren `json:"rens"`
}

// Init - 初期化
func (db *RenDb) Init(boardWidth int, boardHeight int) {
	db.Header.Init(boardWidth, boardHeight)

	// Clear
	for ri := range db.Rens {
		delete(db.Rens, ri)
	}
}

// Save - 連データベースの外部ファイル書込
func (db *RenDb) Save(path string, convertLocation func(point.Point) string, onError func(error) bool) bool {

	// 外部ファイルに出力するための、内部状態の整形
	db.RefreshToExternalFile(convertLocation)

	// Marshal関数でjsonエンコード
	// ->返り値jsonDataにはエンコード結果が[]byteの形で格納される
	jsonBinary, errA := json.Marshal(db)
	if errA != nil {
		return onError(errA)
	}

	// ファイル読込
	var errB = os.WriteFile(path, jsonBinary, 0664)
	if errB != nil {
		return onError(errB)
	}

	return true
}

// FindRen - 連を取得
func (db *RenDb) GetRen(renId ren_id.RenId) (*rentype.Ren, bool) {
	var ren1, isOk = db.Rens[renId]

	if isOk {
		return ren1, true
	}

	return nil, false
}

// RegisterRen - 連を登録
// * すでに Id が登録されているなら、上書きしない
func (db *RenDb) RegisterRen(positionNthFigure int, movesNum1 moves_num.MovesNum, ren1 *rentype.Ren) {
	var renId = GetRenId(db.Header.GetBoardMemoryWidth(), positionNthFigure, movesNum1, ren1.MinimumLocation)

	var _, isExists = db.Rens[renId]
	if !isExists {
		db.Rens[renId] = ren1
	}
}

// Dump - ダンプ
func (db *RenDb) Dump() string {
	var sb strings.Builder

	// 全ての要素
	for idStr, ren1 := range db.Rens {
		sb.WriteString(fmt.Sprintf("[%s]%s \n", idStr, ren1.Dump()))
	}

	var text = sb.String()
	if 0 < len(text) {
		text = text[:len(text)-1]
	}
	return text
}

// RefreshToExternalFile - 外部ファイルに出力されてもいいように内部状態を整形します
func (db *RenDb) RefreshToExternalFile(convertLocation func(point.Point) string) {
	for _, ren1 := range db.Rens {
		ren1.RefreshToExternalFile(convertLocation)
	}
}

// GetRenId - 連のIdを取得
func GetRenId(boardMemoryWidth int, positionNthFigure int, movesNum1 moves_num.MovesNum, minimumLocation point.Point) ren_id.RenId {
	var posNth = movesNum1 + geta.Geta
	var coord = board_coordinate.GetRenIdFromPointOnBoard(boardMemoryWidth, minimumLocation)

	return ren_id.RenId(fmt.Sprintf("%0*d,%s", positionNthFigure, posNth, coord))
}

// NewRenDb - 連データベースを新規作成
func NewRenDb(boardWidth int, boardHeight int) *RenDb {
	var db = new(RenDb)
	db.Header.Init(boardWidth, boardHeight)
	db.Rens = make(map[ren_id.RenId]*rentype.Ren)
	return db
}
