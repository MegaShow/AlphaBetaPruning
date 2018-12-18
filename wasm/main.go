package main

import (
	"math"
	"strings"
	"syscall/js"
)

func registerCallbacks() {
	global := js.Global()
	getAlphaBetaMove := js.NewCallback(func(args []js.Value) {
		strArgs := strings.Split(args[0].String(), " ")
		fen := strings.Split(strArgs[0], "/")
		chess := Chess{Next: byte(strArgs[1][0])}
		for i := 0; i < len(fen); i++ {
			var index int
			for _, v := range fen[len(fen)-1-i] {
				if v >= '1' && v <= '9' {
					index += int(v - '0')
				} else {
					chess.Data[i][index] = byte(v)
					index++
				}
			}
		}
		if len(args) > 1 && args[1].String() != "" {
			steps := strings.Split(args[1].String(), ",")
			for _, step := range steps {
				chess.Data[step[1]-'0'][step[0]-'a'], chess.Data[step[3]-'0'][step[2]-'a'] = 0, chess.Data[step[1]-'0'][step[0]-'a']
				if chess.Next == 'w' {
					chess.Next = 'b'
				} else {
					chess.Next = 'w'
				}
			}
		}
		_, s := AlphaBetaPruning(&chess, 4, math.MinInt32, math.MaxInt32)
		step := string([]rune{rune(s[0] + 'a'), rune(s[1] + '0'), rune(s[2] + 'a'), rune(s[3] + '0')})
		// println(v, step)
		js.Global().Call("move", step)
	})
	global.Set("getAlphaBetaMove", getAlphaBetaMove)
}

func main() {
	c := make(chan struct{}, 0)
	println("hello wasm")
	registerCallbacks()
	<-c
}
