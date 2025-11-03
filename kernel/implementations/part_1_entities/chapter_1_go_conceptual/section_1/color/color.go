package color

import "fmt"

type Color uint

const (
	None Color = iota
	Black
	White
)

// String - 文字列化
func (c Color) String() string {
	switch c {
	case None:
		return ""
	case Black:
		return "x"
	case White:
		return "o"
	default:
		panic(fmt.Sprintf("unexpected color:%d", int(c)))
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
