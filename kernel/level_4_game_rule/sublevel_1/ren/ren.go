package ren

import (
	"fmt"
	"math"
	"strings"

	// Level 2.1
	color "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/color"
	point "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/point"

	// Level 2
	stone "github.com/muzudho/kifuwarabe-uec17/kernel/level_3_physical/sublevel_1/stone"
)

// Ren - 連，れん
type Ren struct {
	// Sto - （外部ファイル向け）石
	Sto string `json:"stone"`
	// Loc - （外部ファイル向け）石の盤上の座標符号の空白区切りのリスト
	Loc string `json:"locate"`
	// LibLoc - （外部ファイル向け）呼吸点の盤上の座標符号の空白区切りのリスト
	LibLoc string `json:"liberty"`

	// 隣接する石の色
	AdjacentColor color.Color
	// 石
	Stone stone.Stone
	// 要素の石の位置
	Locations []point.Point
	// 呼吸点の位置
	LibertyLocations []point.Point
	// 最小の場所。Idとして利用することを想定
	MinimumLocation point.Point
}

// NewRen - 連を新規作成
//
// Parameters
// ----------
// color - 色
func NewRen(stone stone.Stone) *Ren {
	var r = new(Ren)
	r.Stone = stone
	r.AdjacentColor = color.Color_None
	r.MinimumLocation = math.MaxInt
	return r
}

// GetArea - 面積。アゲハマの数
func (r *Ren) GetArea() int {
	return len(r.Locations)
}

// GetLibertyArea - 呼吸点の面積
func (r *Ren) GetLibertyArea() int {
	return len(r.LibertyLocations)
}

// GetStone - 石
func (r *Ren) GetStone() stone.Stone {
	return r.Stone
}

// GetAdjacentColor - 隣接する石の色
func (r *Ren) GetAdjacentColor() color.Color {
	return r.AdjacentColor
}

// GetMinimumLocation - 最小の場所。Idとして利用することを想定
func (r *Ren) GetMinimumLocation() point.Point {
	return r.MinimumLocation
}

// AddLocation - 場所の追加
func (r *Ren) AddLocation(location point.Point) {
	r.Locations = append(r.Locations, location)

	// 最小の数を更新
	r.MinimumLocation = point.Point(math.Min(float64(r.MinimumLocation), float64(location)))
}

// ForeachLocation - 場所毎に
func (r *Ren) ForeachLocation(setLocation func(int, point.Point)) {
	for i, point := range r.Locations {
		setLocation(i, point)
	}
}

// Dump - ダンプ
//
// Example: `22 23 24 25`
func (r *Ren) Dump() string {
	var convertLocation = func(location point.Point) string {
		return fmt.Sprintf("%d", location)
	}
	var tokens = r.createCoordBelt(r.Locations, convertLocation)
	return strings.Join(tokens, " ")
}

// 文字列の配列を作成
// Example: {`22`, `23` `24`, `25`}
func (r *Ren) createCoordBelt(locations []point.Point, convertLocation func(point.Point) string) []string {
	var tokens []string

	// 全ての要素
	for _, location := range locations {
		var token = convertLocation(location)
		tokens = append(tokens, token)
	}

	return tokens
}

// RefreshToExternalFile - 外部ファイルに出力されてもいいように内部状態を整形します
func (r *Ren) RefreshToExternalFile(convertLocation func(point.Point) string) {
	{
		// stone to Sto
		// Examples: `.`, `x`, `o`, `+`
		r.Sto = r.Stone.String()
	}
	{
		// lorations to Loc
		// Example: `A1 B2 C3 D4`
		var tokens = r.createCoordBelt(r.Locations, convertLocation)
		// sort.Strings(tokens) // 辞書順ソート - 走査方向が変わってしまうので止めた
		r.Loc = strings.Join(tokens, " ")
	}
	{
		// libertyLocations to LibLoc
		var tokens = r.createCoordBelt(r.LibertyLocations, convertLocation)
		r.LibLoc = strings.Join(tokens, " ")
	}
}
