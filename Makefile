.PHONY: build

build:
	go generate
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o cooking.linux
	podman build . -t cooking
	podman save cooking > cooking.img

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
