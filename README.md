# event driven web shop
aka Event Thingy Store

## setup
Terminal 1
```
docker compose up
```

Terminal 2
```
cd backend
make start
```

Terminal 3
```
cd frontend
make start
```

## Required tools for development
- golang
- air https://github.com/cosmtrek/air
- elm
- elm-live
- protoc
- protoc-gen-elm https://www.npmjs.com/package/protoc-gen-elm
- docker


## Features
- automatic rebuild backend
- automatic rebuild frontend
- wait free reads
- removed rw lock
- load all data before starting to serve requests

## Rules
- keep read short
- all state exists as persisted event / kafka is the only source of truth
- avoid (syncronous) network calls
- minimize overhead in data flow, data access, code structure

## Roadmap

sync frontend
- websocket
- localstorage
- tab open

backend
- remove .uuid
- after dirty shutdown wait for ip ttl
- use context specific data structs
- snapshots
- share events between contexts
- generalize context
- tests for handlers * cqrs & context
- context with swapable model & read-write lock
- liveness probe / readiness probe
- /metrics endpoint for prometheus
- tracing via jaeger
- add timeouts & retries & exponential backoff & shedding & circuitbreaking & avoid thundering heard
- ensure ordered cart is exactly the cart shown in the browser
- when triggering an sideeffect, then ensure to do this only once
- sync to other contexts before processing commands

frontend
- remove .uuid
- use int32 everywhere https://package.elm-lang.org/packages/eriktim/elm-protocol-buffers/latest/#known-limitations
- use less & cleanup elm
- use protobuf / remove json
- use modules
- use urls & links
- only process results matching the current model
- configure backend url
- tests
- add timeouts & retries & exponential backoff & shedding & circuitbreaking & avoid thundering heard

general
- e2e tests
- tilt
- event browser
- add pim
- add search
- add fulfilment

## Failed ideas
- use arc cache to skip struct to json marshal
