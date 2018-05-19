FROM alpine:3.7
RUN apk -U add ca-certificates
COPY restify /usr/bin/
ENTRYPOINT ["/usr/bin/restify"]