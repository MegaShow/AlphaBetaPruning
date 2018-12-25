package main

import (
	"strings"
	"syscall/js"

	"github.com/MegaShow/AlphaBetaPruning/chess"
)

func registerCallbacks() {
	global := js.Global()
	getAlphaBetaMove := js.NewCallback(func(args []js.Value) {
		fen := args[0].String()
		var steps []string
		if len(args) > 1 {
			steps = strings.Split(args[1].String(), ",")
		}
		step := chess.GetBestMove(fen, steps)
		js.Global().Call("move", step)
	})
	global.Set("getAlphaBetaMoveByWasm", getAlphaBetaMove)
}

func main() {
	c := make(chan struct{}, 0)
	println("hello wasm")
	registerCallbacks()
	<-c
}
