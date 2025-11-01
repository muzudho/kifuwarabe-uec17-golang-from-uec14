// BOF [O11o_4o2o1o0]

package kernel

import (
	"fmt"
	"math"
	"strings"

	types "github.com/muzudho/kifuwarabe-uec17/kernel/types"
)

// Ren - 連，れん
type Ren struct {
	// Sto - （外部ファイル向け）石
	Sto string `json:"stone"`
	// Loc - （外部ファイル向け）石の盤上の座標符号の空白区切りのリスト
	Loc string `json:"locate"`
	// LibLoc - （外部ファイル向け）呼吸点の盤上の座標符号の空白区切りのリスト
	LibLoc string `json:"liberty"`

	// 石
	stone Stone
	// 隣接する石の色
	adjacentColor types.Color
	// 要素の石の位置
	locations []Point
	// 呼吸点の位置
	libertyLocations []Point
	// 最小の場所。Idとして利用することを想定
	minimumLocation Point
}

// NewRen - 連を新規作成
//
// Parameters
// ----------
// color - 色
func NewRen(stone Stone) *Ren {
	var r = new(Ren)
	r.stone = stone
	r.adjacentColor = types.Color_None
	r.minimumLocation = math.MaxInt
	return r
}

// GetArea - 面積。アゲハマの数
func (r *Ren) GetArea() int {
	return len(r.locations)
}

// GetLibertyArea - 呼吸点の面積
func (r *Ren) GetLibertyArea() int {
	return len(r.libertyLocations)
}

// GetStone - 石
func (r *Ren) GetStone() Stone {
	return r.stone
}

// GetAdjacentColor - 隣接する石の色
func (r *Ren) GetAdjacentColor() types.Color {
	return r.adjacentColor
}

// GetMinimumLocation - 最小の場所。Idとして利用することを想定
func (r *Ren) GetMinimumLocation() Point {
	return r.minimumLocation
}

// AddLocation - 場所の追加
func (r *Ren) AddLocation(location Point) {
	r.locations = append(r.locations, location)

	// 最小の数を更新
	r.minimumLocation = Point(math.Min(float64(r.minimumLocation), float64(location)))
}

// ForeachLocation - 場所毎に
func (r *Ren) ForeachLocation(setLocation func(int, Point)) {
	for i, point := range r.locations {
		setLocation(i, point)
	}
}

// Dump - ダンプ
//
// Example: `22 23 24 25`
func (r *Ren) Dump() string {
	var convertLocation = func(location Point) string {
		return fmt.Sprintf("%d", location)
	}
	var tokens = r.createCoordBelt(r.locations, convertLocation)
	return strings.Join(tokens, " ")
}

// 文字列の配列を作成
// Example: {`22`, `23` `24`, `25`}
func (r *Ren) createCoordBelt(locations []Point, convertLocation func(Point) string) []string {
	var tokens []string

	// 全ての要素
	for _, location := range locations {
		var token = convertLocation(location)
		tokens = append(tokens, token)
	}

	return tokens
}

// RefreshToExternalFile - 外部ファイルに出力されてもいいように内部状態を整形します
func (r *Ren) RefreshToExternalFile(convertLocation func(Point) string) {
	{
		// stone to Sto
		// Examples: `.`, `x`, `o`, `+`
		r.Sto = r.stone.String()
	}
	{
		// lorations to Loc
		// Example: `A1 B2 C3 D4`
		var tokens = r.createCoordBelt(r.locations, convertLocation)
		// sort.Strings(tokens) // 辞書順ソート - 走査方向が変わってしまうので止めた
		r.Loc = strings.Join(tokens, " ")
	}
	{
		// libertyLocations to LibLoc
		var tokens = r.createCoordBelt(r.libertyLocations, convertLocation)
		r.LibLoc = strings.Join(tokens, " ")
	}
}

// EOF [O11o_4o2o1o0]
