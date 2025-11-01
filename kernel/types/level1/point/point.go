// BOF [O12o__10o1o0]

package point

// Point - 交点の座標。いわゆる配列のインデックス。枠を含む盤の左上を 0 とします
type Point int

// Cell_Pass - パス
const Cell_Pass Point = 0

// IllegalZ - 石が置けない番地の目印として使用。例：UCT計算中に石が置けなかった
const Cell_Illegal Point = -1

// EOF [O12o__10o1o0]
