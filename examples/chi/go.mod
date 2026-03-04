module github.com/Just-maple/xmux/examples/chi

go 1.18

require (
	github.com/Just-maple/xmux v1.0.0
	github.com/Just-maple/xmux/examples/common v0.0.0
	github.com/go-chi/chi/v5 v5.0.11
)

replace github.com/Just-maple/xmux => ../../

replace github.com/Just-maple/xmux/examples/common => ../common
