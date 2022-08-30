ssh my <<EOF
    docker run \
    --network chat \
    --detach \
    --name cooking \
    --restart always \
    cooking
EOF
