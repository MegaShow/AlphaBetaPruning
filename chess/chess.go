package chess

import (
	"unicode"
)

// Chess 棋盘信息
type Chess struct {
	Data [10][9]byte
	Next byte
}

// 子力价值
var BasicValueTable = map[byte]int{
	// 小写 - 黑, 大写 - 红
	// k - 将, a - 士, b - 象, n - 马, r - 车, c - 炮, p - 兵
	'k': 10000, 'a': 110, 'b': 110, 'n': 300, 'r': 600, 'c': 300, 'p': 70,
}

// 机动性价值
var MobilityValueTable = map[byte]int{
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

func evaluate(chess *Chess) int {
	var redValue, blackValue int
	var redBasicValue, blackBasicValue int
	var redPositionValue, blackPositionValue int
	for n := 0; n < len(chess.Data); n++ {
		for a := 0; a < len(chess.Data[n]); a++ {
			positions := PositionValueTable[byte(unicode.ToLower(rune(chess.Data[n][a])))]
			if chess.Next == 'b' {
				blackBasicValue += BasicValueTable[chess.Data[n][a]]
				blackPositionValue += positions[n][a]
			} else {
				redBasicValue += BasicValueTable[byte(unicode.ToLower(rune(chess.Data[n][a])))]
				redPositionValue += positions[9-n][a]
			}
		}
	}
	redValue = redBasicValue + redPositionValue
	blackValue = blackBasicValue + blackBasicValue
	if chess.Next == 'b' {
		return blackValue - redValue
	} else {
		return redValue - blackValue
	}
}

func checkStep(chess *Chess, steps *[][4]int, oldA, oldN, newA, newN int) {
	if newN < len(chess.Data) && newA < len(chess.Data[newN]) &&
		(chess.Data[newN][newA] == 0 || (chess.Next == 'b' && unicode.IsUpper(rune(chess.Data[newN][newA]))) ||
			(chess.Next != 'b' && unicode.IsLower(rune(chess.Data[newN][newA])))) {
		*steps = append(*steps, [4]int{oldA, oldN, newA, newN})
	}
}

func createNextStep(chess *Chess) (steps [][4]int) {
	if chess.Next == 'b' {
		for n := 0; n < len(chess.Data); n++ {
			for a := 0; a < len(chess.Data[n]); a++ {
				if chess.Data[n][a] == 'k' { // 将
					if n >= 8 {
						checkStep(chess, &steps, a, n, a, n-1)
					}
					if n <= 8 {
						checkStep(chess, &steps, a, n, a, n+1)
					}
					if a >= 'e'-'a' {
						checkStep(chess, &steps, a, n, a-1, n)
					}
					if a <= 'e'-'a' {
						checkStep(chess, &steps, a, n, a+1, n)
					}
				} else if chess.Data[n][a] == 'a' { // 士
					if n == 8 {
						checkStep(chess, &steps, a, n, 'd'-'a', 7)
						checkStep(chess, &steps, a, n, 'd'-'a', 9)
						checkStep(chess, &steps, a, n, 'f'-'a', 7)
						checkStep(chess, &steps, a, n, 'f'-'a', 9)
					} else {
						checkStep(chess, &steps, a, n, 'e'-'a', 8)
					}
				} else if chess.Data[n][a] == 'b' { // 象
					if a == 'a'-'a' {
						checkStep(chess, &steps, a, n, 'c'-'a', 5)
						checkStep(chess, &steps, a, n, 'c'-'a', 9)
					} else if a == 'c'-'a' {
						checkStep(chess, &steps, a, n, 'a'-'a', 7)
						checkStep(chess, &steps, a, n, 'e'-'a', 7)
					} else if a == 'e'-'a' {
						checkStep(chess, &steps, a, n, 'c'-'a', 5)
						checkStep(chess, &steps, a, n, 'c'-'a', 9)
						checkStep(chess, &steps, a, n, 'g'-'a', 5)
						checkStep(chess, &steps, a, n, 'g'-'a', 9)
					} else if a == 'g'-'a' {
						checkStep(chess, &steps, a, n, 'e'-'a', 7)
						checkStep(chess, &steps, a, n, 'i'-'a', 7)
					} else if a == 'i'-'a' {
						checkStep(chess, &steps, a, n, 'g'-'a', 5)
						checkStep(chess, &steps, a, n, 'g'-'a', 9)
					}
				} else if chess.Data[n][a] == 'n' {
					if n+2 < len(chess.Data) && a+1 < len(chess.Data[n]) && chess.Data[n+1][a] == 0 {
						checkStep(chess, &steps, a, n, a+1, n+2)
					}
					if n+1 < len(chess.Data) && a+2 < len(chess.Data[n]) && chess.Data[n][a+1] == 0 {
						checkStep(chess, &steps, a, n, a+2, n+1)
					}
					if n-1 >= 0 && a+2 < len(chess.Data[n]) && chess.Data[n][a+1] == 0 {
						checkStep(chess, &steps, a, n, a+2, n-1)
					}
					if n-2 >= 0 && a+1 < len(chess.Data[n]) && chess.Data[n-1][a] == 0 {
						checkStep(chess, &steps, a, n, a+1, n-2)
					}
					if n-2 >= 0 && a-1 >= 0 && chess.Data[n-1][a] == 0 {
						checkStep(chess, &steps, a, n, a-1, n-2)
					}
					if n-1 >= 0 && a-2 >= 0 && chess.Data[n][a-1] == 0 {
						checkStep(chess, &steps, a, n, a-2, n-1)
					}
					if n+1 < len(chess.Data) && a-2 >= 0 && chess.Data[n][a-1] == 0 {
						checkStep(chess, &steps, a, n, a+2, n-1)
					}
					if n+2 < len(chess.Data) && a-1 >= 0 && chess.Data[n+1][a] == 0 {
						checkStep(chess, &steps, a, n, a-1, n+2)
					}
				} else if chess.Data[n][a] == 'r' { // 车
					for i := 1; n+i < len(chess.Data); i++ {
						checkStep(chess, &steps, a, n, a, n+i)
						if chess.Data[n+i][a] != 0 {
							break
						}
					}
					for i := 1; n-i >= 0; i++ {
						checkStep(chess, &steps, a, n, a, n-i)
						if chess.Data[n-i][a] != 0 {
							break
						}
					}
					for i := 1; a+i < len(chess.Data[n]); i++ {
						checkStep(chess, &steps, a, n, a+i, n)
						if chess.Data[n][a+i] != 0 {
							break
						}
					}
					for i := 1; a-i >= 0; i++ {
						checkStep(chess, &steps, a, n, a-i, n)
						if chess.Data[n][a-i] != 0 {
							break
						}
					}
				} else if chess.Data[n][a] == 'c' { // 炮
					var i int
					for i = 1; n+i < len(chess.Data) && chess.Data[n+i][a] == 0; i++ {
						checkStep(chess, &steps, a, n, a, n+i)
					}
					for i++; n+i < len(chess.Data); i++ {
						if chess.Data[n+i][a] != 0 {
							checkStep(chess, &steps, a, n, a, n+i)
							break
						}
					}
					for i = 1; n-i >= 0 && chess.Data[n-i][a] == 0; i++ {
						checkStep(chess, &steps, a, n, a, n-i)
					}
					for i++; n-i >= 0; i++ {
						if chess.Data[n-i][a] != 0 {
							checkStep(chess, &steps, a, n, a, n-i)
							break
						}
					}
					for i = 1; a+i < len(chess.Data[n]) && chess.Data[n][a+i] == 0; i++ {
						checkStep(chess, &steps, a, n, a+i, n)
					}
					for i++; a+i < len(chess.Data[n]); i++ {
						if chess.Data[n][a+i] != 0 {
							checkStep(chess, &steps, a, n, a+i, n)
							break
						}
					}
					for i = 1; a-i >= 0 && chess.Data[n][a-i] == 0; i++ {
						checkStep(chess, &steps, a, n, a-i, n)
					}
					for i++; a-i >= 0; i++ {
						if chess.Data[n][a-i] != 0 {
							checkStep(chess, &steps, a, n, a-i, n)
							break
						}
					}
				} else if chess.Data[n][a] == 'p' { // 兵
					if n == 6 || n == 5 {
						checkStep(chess, &steps, a, n, a, n-1)
					} else {
						checkStep(chess, &steps, a, n, a, n-1)
						checkStep(chess, &steps, a, n, a-1, n)
						checkStep(chess, &steps, a, n, a+1, n)
					}
				}
			}
		}
	} else {
		for n := 0; n < len(chess.Data); n++ {
			for a := 0; a < len(chess.Data[n]); a++ {
				if chess.Data[n][a] == 'K' { // 帅
					if n >= 1 {
						checkStep(chess, &steps, a, n, a, n-1)
					}
					if n <= 1 {
						checkStep(chess, &steps, a, n, a, n+1)
					}
					if a >= 'e'-'a' {
						checkStep(chess, &steps, a, n, a-1, n)
					}
					if a <= 'e'-'a' {
						checkStep(chess, &steps, a, n, a+1, n)
					}
				} else if chess.Data[n][a] == 'A' { // 仕
					if n == 1 {
						checkStep(chess, &steps, a, n, 'd'-'a', 0)
						checkStep(chess, &steps, a, n, 'd'-'a', 2)
						checkStep(chess, &steps, a, n, 'f'-'a', 0)
						checkStep(chess, &steps, a, n, 'f'-'a', 2)
					} else {
						checkStep(chess, &steps, a, n, 'e'-'a', 1)
					}
				} else if chess.Data[n][a] == 'B' { // 相
					if a == 'a'-'a' {
						checkStep(chess, &steps, a, n, 'c'-'a', 0)
						checkStep(chess, &steps, a, n, 'c'-'a', 4)
					} else if a == 'c'-'a' {
						checkStep(chess, &steps, a, n, 'a'-'a', 2)
						checkStep(chess, &steps, a, n, 'e'-'a', 2)
					} else if a == 'e'-'a' {
						checkStep(chess, &steps, a, n, 'c'-'a', 0)
						checkStep(chess, &steps, a, n, 'c'-'a', 4)
						checkStep(chess, &steps, a, n, 'g'-'a', 0)
						checkStep(chess, &steps, a, n, 'g'-'a', 4)
					} else if a == 'g'-'a' {
						checkStep(chess, &steps, a, n, 'e'-'a', 2)
						checkStep(chess, &steps, a, n, 'i'-'a', 2)
					} else if a == 'i'-'a' {
						checkStep(chess, &steps, a, n, 'g'-'a', 0)
						checkStep(chess, &steps, a, n, 'g'-'a', 4)
					}
				} else if chess.Data[n][a] == 'N' {
					if n+2 < len(chess.Data) && a+1 < len(chess.Data[n]) && chess.Data[n+1][a] == 0 {
						checkStep(chess, &steps, a, n, a+1, n+2)
					}
					if n+1 < len(chess.Data) && a+2 < len(chess.Data[n]) && chess.Data[n][a+1] == 0 {
						checkStep(chess, &steps, a, n, a+2, n+1)
					}
					if n-1 >= 0 && a+2 < len(chess.Data[n]) && chess.Data[n][a+1] == 0 {
						checkStep(chess, &steps, a, n, a+2, n-1)
					}
					if n-2 >= 0 && a+1 < len(chess.Data[n]) && chess.Data[n-1][a] == 0 {
						checkStep(chess, &steps, a, n, a+1, n-2)
					}
					if n-2 >= 0 && a-1 >= 0 && chess.Data[n-1][a] == 0 {
						checkStep(chess, &steps, a, n, a-1, n-2)
					}
					if n-1 >= 0 && a-2 >= 0 && chess.Data[n][a-1] == 0 {
						checkStep(chess, &steps, a, n, a-2, n-1)
					}
					if n+1 < len(chess.Data) && a-2 >= 0 && chess.Data[n][a-1] == 0 {
						checkStep(chess, &steps, a, n, a-2, n+1)
					}
					if n+2 < len(chess.Data) && a-1 >= 0 && chess.Data[n+1][a] == 0 {
						checkStep(chess, &steps, a, n, a-1, n+2)
					}
				} else if chess.Data[n][a] == 'R' { // 车
					for i := 1; n+i < len(chess.Data); i++ {
						checkStep(chess, &steps, a, n, a, n+i)
						if chess.Data[n+i][a] != 0 {
							break
						}
					}
					for i := 1; n-i >= 0; i++ {
						checkStep(chess, &steps, a, n, a, n-i)
						if chess.Data[n-i][a] != 0 {
							break
						}
					}
					for i := 1; a+i < len(chess.Data[n]); i++ {
						checkStep(chess, &steps, a, n, a+i, n)
						if chess.Data[n][a+i] != 0 {
							break
						}
					}
					for i := 1; a-i >= 0; i++ {
						checkStep(chess, &steps, a, n, a-i, n)
						if chess.Data[n][a-i] != 0 {
							break
						}
					}
				} else if chess.Data[n][a] == 'C' { // 炮
					var i int
					for i = 1; n+i < len(chess.Data) && chess.Data[n+i][a] == 0; i++ {
						checkStep(chess, &steps, a, n, a, n+i)
					}
					for i++; n+i < len(chess.Data); i++ {
						if chess.Data[n+i][a] != 0 {
							checkStep(chess, &steps, a, n, a, n+i)
							break
						}
					}
					for i = 1; n-i >= 0 && chess.Data[n-i][a] == 0; i++ {
						checkStep(chess, &steps, a, n, a, n-i)
					}
					for i++; n-i >= 0; i++ {
						if chess.Data[n-i][a] != 0 {
							checkStep(chess, &steps, a, n, a, n-i)
							break
						}
					}
					for i = 1; a+i < len(chess.Data[n]) && chess.Data[n][a+i] == 0; i++ {
						checkStep(chess, &steps, a, n, a+i, n)
					}
					for i++; a+i < len(chess.Data[n]); i++ {
						if chess.Data[n][a+i] != 0 {
							checkStep(chess, &steps, a, n, a+i, n)
							break
						}
					}
					for i = 1; a-i >= 0 && chess.Data[n][a-i] == 0; i++ {
						checkStep(chess, &steps, a, n, a-i, n)
					}
					for i++; a-i >= 0; i++ {
						if chess.Data[n][a-i] != 0 {
							checkStep(chess, &steps, a, n, a-i, n)
							break
						}
					}
				} else if chess.Data[n][a] == 'P' { // 兵
					if n == 3 || n == 4 {
						checkStep(chess, &steps, a, n, a, n+1)
					} else {
						if n+1 < len(chess.Data) {
							checkStep(chess, &steps, a, n, a, n+1)
						}
						if a-1 >= 0 {
							checkStep(chess, &steps, a, n, a-1, n)
						}
						if a+1 < len(chess.Data[n]) {
							checkStep(chess, &steps, a, n, a+1, n)
						}
					}
				}
			}
		}
	}
	return
}

// AlphaBetaPruning alpha-beta剪枝
func AlphaBetaPruning(chess *Chess, depth, alpha, beta int) (int, [4]int) {
	if depth == 0 {
		return evaluate(chess), [4]int{}
	}
	steps := createNextStep(chess)
	// fmt.Println(chess.Next, depth, alpha, beta)
	var old byte
	var s [4]int
	for _, step := range steps {
		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]], old = 0, chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]]
		chess.Next = 'w'
		v, _ := AlphaBetaPruning(chess, depth-1, -beta, -alpha)
		v = -v
		chess.Next = 'b'
		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]] = chess.Data[step[3]][step[2]], old
		if v >= beta {
			return beta, s
		}
		if v > alpha {
			alpha = v
			s = step
		}
	}
	return alpha, s
	//if chess.Next == 'b' { // 极大值层
	//	for _, step := range steps {
	//		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]], old = 0, chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]]
	//		chess.Next = 'w'
	//		v, _ := AlphaBetaPruning(chess, depth-1, alpha, beta)
	//		chess.Next = 'b'
	//		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]] = chess.Data[step[3]][step[2]], old
	//		if v >= alpha {
	//			alpha = v
	//			s = step
	//		}
	//		// beta剪枝
	//		if alpha >= beta {
	//			break
	//		}
	//	}
	//	return alpha, s
	//} else { // 极小值层
	//	for _, step := range steps {
	//		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]], old = 0, chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]]
	//		chess.Next = 'b'
	//		v, _ := AlphaBetaPruning(chess, depth-1, alpha, beta)
	//		chess.Next = 'w'
	//		chess.Data[step[1]][step[0]], chess.Data[step[3]][step[2]] = chess.Data[step[3]][step[2]], old
	//		beta = int(math.Min(float64(v), float64(beta)))
	//		// alpha剪枝
	//		if alpha >= beta {
	//			break
	//		}
	//	}
	//	return beta, s
	//}
}
