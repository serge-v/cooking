.PHONY: build

build:
	GOARCH=amd64 GOOS=linux go build -o cooking.linux

deploy:
	ppodman build . -t cooking
	podman save cooking | gzip | ssh my 'gunzip | docker load'

debug:
	go build
	./cooking -addr localhost:8093

run-local:
	podman \
		--network chat \
		--detach \
		--name cooking \
		--restart always \
		cooking
