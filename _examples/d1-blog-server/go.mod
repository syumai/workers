module github.com/syumai/workers/examples/d1-blog-server

go 1.19

require (
	github.com/mailru/easyjson v0.7.7
	github.com/syumai/workers v0.9.0
)

replace github.com/syumai/workers => ../../

require github.com/josharian/intern v1.0.0 // indirect
