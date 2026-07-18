#!/bin/bash

HOST=podmanuser@my2

ssh ${HOST} <<EOF
    docker kill cooking
    docker rm cooking
    docker image rm localhost/cooking
EOF

gzip -c cooking.img | ssh ${HOST} 'gunzip | docker load'

ssh ${HOST} <<EOF
    docker run \
    --network proxy \
    -p 8080:80 \
    --detach \
    --name cooking \
    --restart always \
    localhost/cooking
EOF
