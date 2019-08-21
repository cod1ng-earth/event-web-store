FROM golang:1.12.7-alpine3.10

RUN apk add --no-cache \
    bash=~5.0 \
    build-base=~0.5 \
    git=~2.22 \
    make=~4.2 \
    ncurses=~6.1 \
    netcat-openbsd=~1.130
RUN wget -O /bin/air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air && \
    chmod +x /bin/air

ENV GO111MODULE=on
WORKDIR /go/src/app

# These are so-called anonymous volumes.
# The rest of the container's file-system can be read-only
# e.g via Docker Compose's 'read_only: true'
# or Kubernetes' pod security policy 'ReadOnlyRootFilesystem'
RUN mkdir -p /.cache/go-build /go/pkg && \
    chmod 777 /.cache/go-build /go/pkg
VOLUME ["/tmp", "/go/pkg", "/.cache/go-build"]

#USER 65635
CMD [ "make", "run" ]
