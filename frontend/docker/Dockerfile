FROM node:12.8.0-alpine

RUN apk add --no-cache \
    bash=~4.4 \
    make=~4.2 \
    ncurses=~6.1
RUN yarn global add elm elm-live

WORKDIR /app

# These are so-called anonymous volumes.
# The rest of the container's file-system can be read-only
# e.g via Docker Compose's 'read_only: true'
# or Kubernetes' pod security policy 'ReadOnlyRootFilesystem'
RUN mkdir -p /home/node/.elm /app/elm-stuff && \
    chmod 777 /home/node/.elm /app/elm-stuff
VOLUME ["/home/node/.elm", "/app/elm-stuff", "/tmp"]

#USER 65635
CMD [ "make", "run" ]
