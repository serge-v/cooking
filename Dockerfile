FROM scratch
WORKDIR /root/
COPY cooking.linux /usr/local/bin/cooking
COPY cooking.zip /root/
CMD ["cooking"]
EXPOSE 80
