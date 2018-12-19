package chess

import (
	"strings"
	"unicode"
)

// Chess 棋盘信息
type Chess struct {
	Board map[int]Piece
	Next  byte
	Size  int
}

type Piece struct {
	Type        int
	Index       int // Index = NumberIndex * 9 + AlphaIndex
	AlphaIndex  int
	NumberIndex int
}

func NewChess(fen string, steps []string) *Chess {
	strArgs := strings.Split(fen, " ")
	lines := strings.Split(strArgs[0], "/")
	chess := Chess{
		Board: make(map[int]Piece),
		Next:  byte(strArgs[1][0]),
		Size:  0,
	}
	for i := 0; i < len(lines); i++ {
		var index int
		for _, v := range lines[len(lines)-1-i] {
			if v >= '1' && v <= '9' {
				index += int(v - '0')
			} else {
				chess.Board[i*9+index] = Piece{
					Type:        int(v),
					Index:       i*9 + index,
					AlphaIndex:  index,
					NumberIndex: i,
				}
				index++
			}
		}
	}
	for _, step := range steps {
		chess.Move([4]int{int(step[1] - '0'), int(step[0] - 'a'), int(step[3] - '0'), int(step[2] - 'a')})
	}
	return &chess
}

func (c *Chess) GetPiecesSize() int {
	return len(c.Board)
}

func (c *Chess) IsWin(player byte) bool {
	return false
}

func (c *Chess) Evaluate(steps [][4]int) int {
	var redValue, blackValue int
	var redBasicValue, blackBasicValue int
	var redPositionValue, blackPositionValue int
	for _, p := range c.Board {
		positions, ok := PositionValueTable[byte(unicode.ToLower(rune(p.Type)))]
		if p.Type >= 'a' && p.Type <= 'z' {
			blackBasicValue += BasicValueTable[p.Type]
			if ok {
				blackPositionValue += positions[p.NumberIndex][p.AlphaIndex]
			}
		} else {
			redBasicValue += BasicValueTable[p.Type]
			if ok {
				redPositionValue += positions[9-p.NumberIndex][p.AlphaIndex]
			}
		}
	}
	redValue = redBasicValue + redPositionValue*8
	blackValue = blackBasicValue + blackPositionValue*8
	// fmt.Println(redValue, blackValue, redBasicValue, redPositionValue, blackBasicValue, blackPositionValue)
	if c.Next == 'b' {
		return blackValue - redValue
	} else {
		return redValue - blackValue
	}
}

// 子力价值
var BasicValueTable = map[int]int{
	// 小写 - 黑, 大写 - 红
	// k - 将, a - 士, b - 象, n - 马, r - 车, c - 炮, p - 兵
	'k': 100000, 'a': 110, 'b': 110, 'n': 300, 'r': 600, 'c': 300, 'p': 70,
	'K': 100000, 'A': 110, 'B': 110, 'N': 300, 'R': 600, 'C': 300, 'P': 70,
}

// 机动性价值
var MobilityValueTable = map[int]int{
	'k': 0, 'a': 1, 'b': 1, 'n': 13, 'r': 7, 'c': 7, 'p': 15,
	'K': 0, 'A': 1, 'B': 1, 'N': 13, 'R': 7, 'C': 7, 'P': 15,
}

// 控制区域价值
var PositionValueTable = map[byte][10][9]int{
	'p': {
		{0, 3, 6, 9, 12, 9, 6, 3, 0},
		{18, 36, 56, 80, 120, 80, 56, 36, 18},
		{14, 26, 42, 60, 80, 60, 42, 26, 14},
		{10, 20, 30, 34, 40, 34, 30, 20, 10},
		{6, 12, 18, 18, 20, 18, 18, 12, 6},
		{2, 0, 8, 0, 8, 0, 8, 0, 2},
		{0, 0, -2, 0, 4, 0, -2, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	'r': {
		{14, 14, 12, 18, 16, 18, 12, 14, 14},
		{16, 20, 18, 24, 26, 24, 18, 20, 16},
		{12, 12, 12, 18, 18, 18, 12, 12, 12},
		{12, 18, 16, 22, 22, 22, 16, 18, 12},
		{12, 14, 12, 18, 18, 18, 12, 14, 12},
		{12, 16, 14, 20, 20, 20, 14, 16, 12},
		{6, 10, 8, 14, 14, 14, 8, 10, 6},
		{4, 8, 6, 14, 12, 14, 6, 8, 4},
		{8, 4, 8, 16, 8, 16, 8, 4, 8},
		{-1, 10, 6, 14, 12, 14, 6, 10, -2},
	},
	'n': {
		{4, 8, 16, 12, 4, 12, 16, 8, 4},
		{4, 10, 28, 16, 8, 16, 28, 10, 4},
		{12, 14, 16, 20, 18, 20, 16, 14, 12},
		{8, 24, 18, 24, 20, 24, 18, 24, 8},
		{6, 16, 14, 18, 16, 18, 14, 16, 6},
		{4, 12, 16, 14, 12, 14, 16, 12, 4},
		{2, 6, 8, 6, 10, 6, 8, 6, 2},
		{4, 2, 8, 8, 4, 8, 8, 2, 4},
		{0, 2, 4, 4, -2, 4, 4, 2, 0},
		{0, -4, 0, 0, 0, 0, 0, -4, 0},
	},
	'c': {
		{6, 4, 0, -10, -12, -10, 0, 4, 6},
		{2, 2, 0, -4, -14, -4, 0, 2, 2},
		{2, 2, 0, -10, -8, -10, 0, 2, 2},
		{0, 0, -2, 4, 10, 4, -2, 0, 0},
		{0, 0, 0, 2, 8, 2, 0, 0, 0},
		{-2, 0, 4, 2, 6, 2, 4, 0, -2},
		{0, 0, 0, 2, 4, 2, 0, 0, 0},
		{4, 0, 8, 6, 10, 6, 8, 0, 4},
		{0, 2, 4, 6, 6, 6, 4, 2, 0},
		{0, 0, 2, 6, 6, 6, 2, 0, 0},
	},
}
