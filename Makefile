all: website

generator: gen.go
	go build -o generator gen.go

website: generator
	./generator
