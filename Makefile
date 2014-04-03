all: redis-trib
build: redis-trib

deps: Godeps
	gpm

sources=cluster.go main.go node.go utils.go

redis-trib: $(sources)
	go build -o redis-trib $(sources)

run: build
	./redis-trib each 127.0.0.1:7001 get foo
