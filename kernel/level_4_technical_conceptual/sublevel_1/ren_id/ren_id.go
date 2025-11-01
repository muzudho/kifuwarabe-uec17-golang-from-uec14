package ren_id

// RenId - 連データベースに格納される連のId
// - 外部ファイルの可読性を優先して数値型ではなく文字列
// - 昇順に並ぶように前ゼロを付ける
type RenId string
