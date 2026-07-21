.PHONY: build

build:
	go generate
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o build/cooking.linux
	podman build --build-arg APP_NAME=cooking -t cooking .
	podman save cooking > build/cooking.img

debug:
	go build -o build/cooking
	build/cooking -addr localhost:8093

run-local:
	podman \
		--network chat \
		--detach \
		--name cooking \
		--restart always \
		cooking
