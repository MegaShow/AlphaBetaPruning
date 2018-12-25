# Alpha-Beta Pruning

Alpha-beta pruning implementation for Chinese chess.

## Dependency

* Golang 1.11+, with Web Assembly for Golang.
* [vschess](https://github.com/FastLight126/vschess), Chinese Chess Web UI.

## Run with Golang Mux

First, download the sources.

```sh
$ go get github.com/MegaShow/AlphaBetaPruning
```

Then, run and listen at the port `8080`.

```sh
$ go run main.go -listen :8080
```

## Run with Web Assembly

> **Suggest that run with `Golang Mux` for better performance.**

First, download the sources.

```sh
$ go get github.com/MegaShow/AlphaBetaPruning
```

Then, build web assembly binary file with `Powershell`.

```sh
$ .\wasm\build.ps1
```

And you need to modify the source file `public/index.html`.

```js
chess = new vschess.load(".vschess", {
  clickResponse: vschess.code.clickResponse.red,
  afterClickAnimate: function () {
    getAlphaBetaMoveByWasm(this.getCurrentFen())
    // getAlphaBetaMoveByApi(this.getCurrentFen())
  }
})
```

Finally, pack the static folder `public` as website root. Or you can use Golang to open a file server.

```sh
$ go run main.go -listen :8080
```

