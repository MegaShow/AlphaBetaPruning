package chess

import (
	"fmt"
	"math"
	"time"
)

func Search(chess *Chess) string {
	start := time.Now()
	depth := 2
	if chess.GetPiecesSize() < 6 {
		depth = 5
	} else if chess.GetPiecesSize() < 12 {
		depth = 4
	} else if chess.GetPiecesSize() < 28 {
		depth = 3
	}
	steps := chess.GetNextSteps()
	alpha, _ := math.MinInt32, math.MaxInt32
	var s [4]int
	for _, step := range steps {
		move := chess.Move(step)
		v := AlphaBetaPruning(chess, depth, alpha, math.MaxInt32)
		// v := Minimax(chess, depth)
		chess.UnMove(move)
		if alpha < v {
			alpha = v
			s = step
		}
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
