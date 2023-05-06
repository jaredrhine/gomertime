module github.com/jaredrhine/gomertime

go 1.20

require (
	github.com/buger/goterm v1.0.4
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	golang.org/x/exp v0.0.0-20230425010034-47ecfdc1ba53
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/klauspost/compress v1.10.3 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/jaredrhine/gomertime => ./pkg/gomertime
