package chess

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func Search(chess *Chess) string {
	start := time.Now()
	depth := 3
	if chess.GetPiecesSize() < 4 {
		depth = 6
	} else if chess.GetPiecesSize() < 8 {
		depth = 5
	} else if chess.GetPiecesSize() < 16 {
		depth = 4
	}
	steps := chess.GetNextSteps()
	alpha, _ := math.MinInt32, math.MaxInt32
	var s [4]int
	var m [][5]int
	for _, step := range steps {
		move := chess.Move(step)
		v := AlphaBetaPruning(chess, depth, alpha, math.MaxInt32)
		// v := Minimax(chess, depth)
		if v <= -90000 {
			ksteps := chess.GetNextSteps()
			var fail bool
			for _, ks := range ksteps {
				move := chess.Move(ks)
				if chess.IsWin('w') {
					chess.UnMove(move)
					fail = true
					break
				}
				chess.UnMove(move)
			}
			if !fail {
				m = append(m, [5]int{v, step[0], step[1], step[2], step[3]})
			}
		}
		chess.UnMove(move)
		if alpha < v {
			alpha = v
			s = step
		}
	}
	if len(m) != 0 {
		sort.Slice(m, func(i, j int) bool {
			return m[i][0] < m[j][0]
		})
		alpha = m[0][0]
		s = [4]int{m[0][1], m[0][2], m[0][3], m[0][4]}
	}
	step := string([]rune{rune(s[1] + 'a'), rune(s[0] + '0'), rune(s[3] + 'a'), rune(s[2] + '0')})
	fmt.Printf("Value: %d, Step: %s, Pieces: %d, Time: %v\n", alpha, step, chess.GetPiecesSize(), time.Since(start))
	return step
}

// Minimax 极小极大搜索算法
func Minimax(chess *Chess, depth int) int {
	if depth == 0 || chess.IsWin('b') || chess.IsWin('w') {
		return chess.Evaluate()
	}
	steps := chess.GetNextSteps()
	if chess.Next == 'b' { // 极大值层
		v := math.MinInt32
		for _, step := range steps {
			move := chess.Move(step)
			v = int(math.Max(float64(v), float64(Minimax(chess, depth-1))))
			chess.UnMove(move)
		}
		return v
	} else { // 极小值层
		v := math.MaxInt32
		for _, step := range steps {
			move := chess.Move(step)
			v = int(math.Min(float64(v), float64(Minimax(chess, depth-1))))
			chess.UnMove(move)
		}
		return v
	}
}

// AlphaBetaPruning alpha-beta剪枝
func AlphaBetaPruning(chess *Chess, depth, alpha, beta int) int {
	if depth == 0 || chess.IsWin('b') || chess.IsWin('w') {
		return chess.Evaluate()
	}
	steps := chess.GetNextSteps()
	for _, step := range steps {
		move := chess.Move(step)
		v := AlphaBetaPruning(chess, depth-1, alpha, beta)
		chess.UnMove(move)
		if chess.Next == 'b' { // 极大值层
			if alpha < v {
				alpha = v
			}
		} else { // 极小值层
			if beta > v {
				beta = v
			}
		}
		// alpha-beta剪枝
		if beta <= alpha {
			break
		}
	}
	if chess.Next == 'b' {
		return alpha
	} else {
		return beta
	}
}
