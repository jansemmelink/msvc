module github.com/jansemmelink/msvc

go 1.12

replace github.com/jansemmelink/config => ../config

require (
	github.com/jansemmelink/config v0.0.0-00010101000000-000000000000
	github.com/jansemmelink/log v0.3.0
	github.com/nats-io/nats.go v1.8.1
)
