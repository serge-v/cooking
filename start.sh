ssh my <<EOF
    docker kill localhost/cooking
    docker rm localhost/cooking
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
