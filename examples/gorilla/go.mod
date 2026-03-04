module github.com/Just-maple/xmux/examples/gorilla

go 1.18

require (
	github.com/Just-maple/xmux v1.0.0
	github.com/Just-maple/xmux/examples/common v0.0.0
	github.com/gorilla/mux v1.8.1
)

replace github.com/Just-maple/xmux => ../../

replace github.com/Just-maple/xmux/examples/common => ../common
