// BOF [O12o__11o__10o5o__10o0]

package level_31_controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// Entities
	point "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/point"
	stone "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/implementations/part_1_entities/chapter_1_go_conceptual/section_1/stone"

	// Level 4.1
	rentype "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_4_game_rule/sublevel_1/ren"

	// Level 6.3
	ren_db "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_6_database/sublevel_3/ren_db"

	// Level 7.1
	liberty_search_algorithm "github.com/muzudho/kifuwarabe-uec17-golang-from-uec14/kernel/level_7_misc/sublevel_1/liberty_search_algorithm"
)

// LoadRenDb - [O12o__11o__10o5o__10o_10o0] 連データベースの外部ファイル読取
func (kernel1 *Kernel) LoadRenDb(path string, onError func(error) bool) bool {
	// ファイル読込
	var binary, errA = os.ReadFile(path)
	if errA != nil {
		return onError(errA)
	}

	var db = new(ren_db.RenDb)
	var errB = json.Unmarshal(binary, db)
	if errB != nil {
		return onError(errB)
	}

	// 外部ファイルからの入力を、内部状態へ適用
	for _, ren := range db.Rens {
		var isOk = kernel1.RefreshRenToInternal(ren)
		if !isOk {
			return false
		}
	}

	// 差し替え
	kernel1.RenDb = db
	return true
}

// RefreshRenToInternal - TODO 外部ファイルから入力された内容を内部状態に適用します
func (kernel1 *Kernel) RefreshRenToInternal(r *rentype.Ren) bool {
	{
		var getDefaultStone = func() (bool, stone.Stone) {
			panic(fmt.Sprintf("unexpected stone:%s", r.Sto))
		}

		// TODO stone from r.Sto
		// Example: "x" --> black
		var isOk, stone = stone.GetStoneFromChar(r.Sto, getDefaultStone)
		if !isOk {
			return false
		}
		r.Stone = stone
	}
	{
		// TODO locations from r.Loc
		// Example: "C1 D1 E1"
		if 0 < len(r.Loc) {
			var codes = strings.Split(r.Loc, " ")

			var numbers = []point.Point{}
			for _, code := range codes {
				var location = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(code)
				numbers = append(numbers, location)
			}

			r.Locations = numbers
		}
	}
	{
		// TODO libertyLocations from r.LibLoc
		// Example: "F1 E2 D2 B1 C2"
		if 0 < len(r.LibLoc) {
			var codes = strings.Split(r.LibLoc, " ")

			var numbers = []point.Point{}
			for _, code := range codes {
				var location = kernel1.Position.Board.Coordinate.GetPointFromGtpMove(code)
				numbers = append(numbers, location)
			}

			r.LibertyLocations = numbers
		}
	}

	return true
}

// RemoveRen - 石の連を打ち上げます
func (kernel1 *Kernel) RemoveRen(ren *rentype.Ren) {
	// 空点をセット
	var setLocation = func(i int, location point.Point) {
		kernel1.Position.Board.SetStoneAt(location, stone.None)
	}

	// 場所毎に
	ren.ForeachLocation(setLocation)
}

// FindAllRens - [O23o_2o1o0] 盤上の全ての連を見つけます
// * 見つけた連は、連データベースへ入れます
func (kernel1 *Kernel) FindAllRens() {
	// チェックボードの初期化
	kernel1.Position.CheckBoard.Init(kernel1.Position.Board.Coordinate)

	var maxPosNthFigure = kernel1.Record.GetMaxPosNthFigure()

	var setLocation = func(location point.Point) {

		var libertySearchAlgorithm = liberty_search_algorithm.NewLibertySearchAlgorithm(kernel1.Position.Board, kernel1.Position.CheckBoard)
		var ren, isFound = libertySearchAlgorithm.FindRen(location)

		if isFound {
			kernel1.RenDb.RegisterRen(maxPosNthFigure, kernel1.Record.MovesNum1, ren)
		}
	}
	// 盤上の枠の内側をスキャン。筋、段の順
	kernel1.Position.Board.GetCoordinate().ForeachPayloadLocationOrderByYx(setLocation)
}
