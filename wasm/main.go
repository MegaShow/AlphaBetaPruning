package main

import "syscall/js"

func registerCallbacks() {
	global := js.Global()
	getAlphaBetaMove := js.NewCallback(func(args []js.Value) {
		println("test")
	})
	global.Set("getAlphaBetaMove", getAlphaBetaMove)
}

func main() {
	c := make(chan struct{}, 0)
	println("hello wasm")
	registerCallbacks()
	<-c
}
