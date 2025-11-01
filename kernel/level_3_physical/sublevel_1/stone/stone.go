package stone

import (
	"fmt"

	// Level 2.1
	color "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_1/color"
)

// Stone - 石の色
type Stone uint

const (
	// Stone_Space - 空点
	Stone_Space Stone = iota
	// Stone_Black - 黒石
	Stone_Black
	// Stone_White - 白石
	Stone_White
	// Stone_Wall - 枠
	Stone_Wall
)

// GetStoneFromName - 文字列の名前を与えると、Stone値を返します
//
// Returns
// -------
// isOk : bool
// stone : Stone
func GetStoneFromName(stoneName string, getDefaultStone func() (bool, Stone)) (bool, Stone) {
	switch stoneName {
	case "space":
		return true, Stone_Space
	case "black":
		return true, Stone_Black
	case "white":
		return true, Stone_White
	case "wall":
		return true, Stone_Wall
	default:
		return getDefaultStone()
	}
}

// GetStoneOrDefaultFromTurn - black または white を与えると、Stone値を返します
//
// Returns
// -------
// stone : Stone
func GetStoneOrDefaultFromTurn(stoneName string, getDefaultStone func() Stone) Stone {
	switch stoneName {
	case "black":
		return Stone_Black
	case "white":
		return Stone_White
	default:
		return getDefaultStone()
	}
}

// GetStoneFromChar - １文字与えると、Stone値を返します
//
// Returns
// -------
// isOk : bool
// stone : Stone
func GetStoneFromChar(stoneChar string, getDefaultStone func() (bool, Stone)) (bool, Stone) {
	switch stoneChar {
	case ".":
		return true, Stone_Space
	case "x":
		return true, Stone_Black
	case "o":
		return true, Stone_White
	case "+":
		return true, Stone_Wall
	default:
		return getDefaultStone()
	}
}

// String - 文字列化
func (s Stone) String() string {
	switch s {
	case Stone_Space:
		return "."
	case Stone_Black:
		return "x"
	case Stone_White:
		return "o"
	case Stone_Wall:
		return "+"
	default:
		panic(fmt.Sprintf("%d", int(s)))
	}
}

// GetColor - 色の取得
func (s Stone) GetColor() color.Color {
	switch s {
	case Stone_Space:
		return color.Color_None
	case Stone_Black:
		return color.Color_Black
	case Stone_White:
		return color.Color_White
	case Stone_Wall:
		return color.Color_None
	default:
		panic(fmt.Sprintf("%d", int(s)))
	}
}
