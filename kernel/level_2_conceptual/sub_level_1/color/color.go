package color

import "fmt"

type Color uint

const (
	Color_None Color = iota
	Color_Black
	Color_White
)

// String - 文字列化
func (c Color) String() string {
	switch c {
	case Color_None:
		return ""
	case Color_Black:
		return "x"
	case Color_White:
		return "o"
	default:
		panic(fmt.Sprintf("unexpected color:%d", int(c)))
	}
}

// GetAdded - 色の加算。上書きはできない
func (c1 Color) GetAdded(c2 Color) Color {
	switch c1 {
	case Color_None:
		return c2
	case Color_Black:
		switch c2 {
		case Color_None:
			return Color_Black
		case Color_Black:
			return Color_Black
		case Color_White:
			return Color_Black
		default:
			panic(fmt.Sprintf("unexpected my_color:%s adds_color:%s", c1, c2))
		}
	case Color_White:
		switch c2 {
		case Color_None:
			return Color_White
		case Color_Black:
			return Color_White
		case Color_White:
			return Color_White
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
	case Color_None:
		return c
	case Color_Black:
		return Color_White
	case Color_White:
		return Color_Black
	default:
		panic(fmt.Sprintf("unexpected color:%d", int(c)))
	}
}
