package chess

import (
	"fmt"
	"math"
)

func Search(chess *Chess) string {
	depth := 2
	if chess.GetPiecesSize() < 8 {
		depth = 5
	} else if chess.GetPiecesSize() < 16 {
		depth = 4
	} else if chess.GetPiecesSize() < 24 {
		depth = 3
	}
	steps := chess.GetNextSteps()
	alpha, beta := math.MinInt32, math.MaxInt32
	var s [4]int
	for _, step := range steps {
		move := chess.Move(step)
		v := AlphaBetaPruning(chess, depth-1, alpha, beta)
		chess.UnMove(move)
		if alpha <= v {
			alpha = v
			s = step
		}
		fmt.Println("step:", step, "alpha:", v)
	}
	step := string([]rune{rune(s[1] + 'a'), rune(s[0] + '0'), rune(s[3] + 'a'), rune(s[2] + '0')})
	fmt.Println(alpha, step, chess.GetPiecesSize())
	return step
}

// AlphaBetaPruning alpha-beta剪枝
func AlphaBetaPruning(chess *Chess, depth, alpha, beta int) int {
	steps := chess.GetNextSteps()
	if chess.IsWin('b') {
		return math.MaxInt32
	} else if chess.IsWin('w') {
		return math.MinInt32
	} else if depth == 0 {
		return chess.Evaluate(steps)
	}
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
