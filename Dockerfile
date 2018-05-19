FROM scratch
COPY restify /usr/bin
ENTRYPOINT ["/usr/bin/restify"]