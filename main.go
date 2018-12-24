package main

import (
	"flag"
	"github.com/MegaShow/AlphaBetaPruning/chess"
	"log"
	"net/http"
	"strings"
)

var (
	listen = flag.String("listen", ":25566", "listen address")
	dir    = flag.String("dir", "./public", "directory to serve")
)

type Body struct {
	Fen   string
	Steps string
}

func main() {
	flag.Parse()
	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/best" {
			err := req.ParseForm()
			if err != nil {
				return
			}
			step := chess.GetBestMove(req.Form.Get("fen"), strings.Split(req.Form.Get("move"), ","))
			resp.Write([]byte(step))
			return
		}
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		http.FileServer(http.Dir(*dir)).ServeHTTP(resp, req)
	})))
}
