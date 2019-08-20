FROM node:12.8.0-alpine
# These are so-called anonymous volumes.
# The rest of the container's file-system can be read-only
# e.g via Docker Compose's 'read_only: true'
# or Kubernetes' pod security policy 'ReadOnlyRootFilesystem'
VOLUME ["/app/elm-stuff", "/.elm", "/root/.elm", "/tmp"]
RUN set -euxo pipefail; \
    apk add --no-cache \
    bash=~4.4 \
    make=~4.2 \
    ncurses=~6.1
RUN set -euxo pipefail; \
    yarn global add elm elm-live
WORKDIR /app
#USER 65635
CMD [ "make", "run" ]
