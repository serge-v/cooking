all: generator

generator: gen.go
	go build -o generator gen.go

website: generator
	./generator

compare:
	# TODO: copy from webdav
	diff prod~ build

deploy:
	# TODO: copy tar, check hash, switch 
