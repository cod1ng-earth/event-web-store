FROM golang:1.12.7-alpine3.10
# These are so-called anonymous volumes.
# The rest of the container's file-system can be read-only
# e.g via Docker Compose's 'read_only: true'
# or Kubernetes' pod security policy 'ReadOnlyRootFilesystem'
VOLUME ["/go/src/app/tmp", "/tmp", "/root/.cache/go-build"]
RUN set -euxo pipefail; \
    apk add --no-cache \
    bash=~5.0 \
    build-base=~0.5 \
    git=~2.22 \
    make=~4.2 \
    ncurses=~6.1 \
    netcat-openbsd=~1.130
RUN set -euxo pipefail; \
    wget -O /bin/air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air \
    && chmod +x /bin/air
WORKDIR /go/src/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN set -euxo pipefail; \
    GO111MODULE=on go mod download
COPY . .
#USER 65635
CMD [ "make", "start" ]
