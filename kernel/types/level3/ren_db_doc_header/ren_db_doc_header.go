package ren_db_doc_header

// Level 2
import "github.com/muzudho/kifuwarabe-uec17/kernel/level_2_conceptual/sublevel_2/board_coordinate"

// RenDbDocHeader - ヘッダー
type RenDbDocHeader struct {
	// BoardWidth - 盤の横幅
	BoardWidth int `json:"boardWidth"`
	// BoardHeight - 盤の縦幅
	BoardHeight int `json:"boardHeight"`
}

// Init - 初期化
func (h *RenDbDocHeader) Init(boardWidth int, boardHeight int) {
	h.BoardWidth = boardWidth
	h.BoardHeight = boardHeight
}

// GetBoardMemoryArea - 枠付き盤の面積
func (h *RenDbDocHeader) GetBoardMemoryArea() int {
	return h.GetBoardMemoryWidth() * h.GetBoardMemoryHeight()
}

// GetBoardMemoryWidth - 枠付き盤の横幅
func (h *RenDbDocHeader) GetBoardMemoryWidth() int {
	return h.BoardWidth + board_coordinate.BothSidesWallThickness
}

// GetBoardMemoryHeight - 枠付き盤の縦幅
func (h *RenDbDocHeader) GetBoardMemoryHeight() int {
	return h.BoardHeight + board_coordinate.BothSidesWallThickness
}
