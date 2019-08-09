FROM node:12.8.0-alpine
RUN set -euxo pipefail; \
    apk add --no-cache \
    bash=~4.4 \
    make=~4.2 \
    ncurses=~6.1
RUN set -euxo pipefail; \
    yarn global add elm elm-live
WORKDIR /app
COPY . .
CMD [ "make", "start" ]
