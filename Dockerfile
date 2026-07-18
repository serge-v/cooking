FROM alpine:3.20
RUN adduser -D -H -s /sbin/nologin cooking
WORKDIR /home/cooking/
COPY build/cooking.linux /usr/local/bin/cooking
COPY build/cooking.zip /home/cooking/
RUN chown -R cooking:cooking /home/cooking/
USER cooking
CMD ["cooking"]
EXPOSE 80
