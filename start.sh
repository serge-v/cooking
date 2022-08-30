ssh my <<EOF
    docker kill cooking
    docker image rm localhost/cooking
EOF

gzip -c cooking.img | ssh my 'gunzip | docker load'

ssh my <<EOF
    docker run \
    --network chat \
    --detach \
    --name cooking \
    --restart always \
    localhost/cooking
EOF
