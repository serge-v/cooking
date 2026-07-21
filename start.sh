#!/bin/bash

HOST=podmanuser@my2

ssh ${HOST} <<EOF
    podman kill cooking
    podman rm cooking
    podman image rm localhost/cooking
EOF

gzip -c build/cooking.img | ssh ${HOST} 'gunzip | podman load'

ssh ${HOST} <<EOF
	podman network exists cooking-net || podman network create cooking-net

	docker run \
		--network cooking-net \
		--detach \
		--name cooking \
		--restart always \
		localhost/cooking
		
	podman network disconnect cooking-net proxy || true
	podman network connect cooking-net proxy
EOF
