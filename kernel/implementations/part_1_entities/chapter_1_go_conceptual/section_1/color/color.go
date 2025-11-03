package color

import "fmt"

// Color - 石の色
type Color uint

const (
	// None - 空点
	None Color = iota
	// Black - 黒石
	Black
	// White - 白石
	White
	// Wall - 枠
	Wall
)

// String - 文字列化
func (c Color) String() string {
	switch c {
	case None:
		return "."
	case Black:
		return "x"
	case White:
		return "o"
	case Wall:
		return "+"
	default:
		panic(fmt.Sprintf("unexpected color:%d", int(c)))
	}
}

// GetStoneFromChar - １文字与えると、Stone値を返します
//
// Returns
// -------
// isOk : bool
// stone : Stone
func GetColorFromCode(stoneChar string, getDefaultColor func() (bool, Color)) (bool, Color) {
	switch stoneChar {
	case ".":
		return true, None
	case "x":
		return true, Black
	case "o":
		return true, White
	case "+":
		return true, Wall
	default:
		return getDefaultColor()
	}
}

// GetColorFromName - 文字列の名前を与えると、Stone値を返します
//
// Returns
// -------
// isOk : bool
// stone : Stone
func GetColorFromName(colorName string, getDefaultColor func() (bool, Color)) (bool, Color) {
	switch colorName {
	case "space":
		return true, None
	case "black":
		return true, Black
	case "white":
		return true, White
	case "wall":
		return true, Wall
	default:
		return getDefaultColor()
	}
}

// GetStoneOrDefaultFromTurn - black または white を与えると、Stone値を返します
//
// Returns
// -------
// stone : Stone
func GetStoneOrDefaultFromTurn(colorName string, getDefaultStone func() Color) Color {
	switch colorName {
	case "space":
		return None
	case "black":
		return Black
	case "white":
		return White
	case "wall":
		return Wall
	default:
		return getDefaultStone()
	}
}

// GetAdded - 色の加算。上書きはできない
func (c1 Color) GetAdded(c2 Color) Color {
	switch c1 {
	case None:
		return c2
	case Black:
		switch c2 {
		case None:
			return Black
		case Black:
			return Black
		case White:
			return Black
		default:
			panic(fmt.Sprintf("unexpected my_color:%s adds_color:%s", c1, c2))
		}
	case White:
		switch c2 {
		case None:
			return White
		case Black:
			return White
		case White:
			return White
		default:
			panic(fmt.Sprintf("unexpected my_color:%s adds_color:%s", c1, c2))
		}
	default:
		panic(fmt.Sprintf("unexpected my_color:%s adds_color:%s", c1, c2))
	}
}

// GetOpponent - 色の反転
func (c Color) GetOpponent() Color {
	switch c {
	case None:
		return c
	case Black:
		return White
	case White:
		return Black
	default:
		panic(fmt.Sprintf("unexpected color:%d", int(c)))
	}
}
