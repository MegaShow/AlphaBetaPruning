package chess

import (
	"fmt"
	"math"
	"strings"
)

func GetBestMove(fen string, steps []string) string {
	strArgs := strings.Split(fen, " ")
	lines := strings.Split(strArgs[0], "/")
	chess := Chess{Next: byte(strArgs[1][0])}
	for i := 0; i < len(lines); i++ {
		var index int
		for _, v := range lines[len(lines)-1-i] {
			if v >= '1' && v <= '9' {
				index += int(v - '0')
			} else {
				chess.Data[i][index] = byte(v)
				index++
			}
		}
	}
	fmt.Println(steps)
	if len(steps) >= 1 && steps[0] != "" {
		for _, step := range steps {
			chess.Data[step[1]-'0'][step[0]-'a'], chess.Data[step[3]-'0'][step[2]-'a'] = 0, chess.Data[step[1]-'0'][step[0]-'a']
			if chess.Next == 'w' {
				chess.Next = 'b'
			} else {
				chess.Next = 'w'
			}
		}
	}
	for n := len(chess.Data) - 1; n >= 0; n-- {
		for a := 0; a < len(chess.Data[n]); a++ {
			if chess.Data[n][a] == 0 {
				fmt.Print("-")
			} else {
				fmt.Print(string(rune(chess.Data[n][a])))
			}
		}
		fmt.Println()
	}
	fmt.Println()
	_, s := AlphaBetaPruning(&chess, 4, math.MinInt32, math.MaxInt32)
	step := string([]rune{rune(s[0] + 'a'), rune(s[1] + '0'), rune(s[2] + 'a'), rune(s[3] + '0')})
	return step
}
