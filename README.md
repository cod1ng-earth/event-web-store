# event driven web shop

aka Event Thingy Store

## setup in docker

``` bash
make setup
```

## setup on host

Terminal 1

``` bash
cd shared
make run
```

Terminal 2

``` bash
cd backend
make import
make run
```

Terminal 3

``` bash
cd frontend
make run
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

- extract simba
  - [DONE] lock: exclusive, sharedread, waitfree
  - [SKIP] serialize: json, protobuf
  - [DONE] configure: massage wrapper
  - [DONE] processing: batch & finalize vs single-step
  - batch only during start vs dynamic selection
  - [DONE] sync to other contexts before processing commands
- extend simba to create bridges
  - share events between contexts
- CORS
  - configurable api domain in frontend
  - configurable spa domain in backend
  - http handler wrapper for CORS setup
- production docker image
- limit cpu and memory use locally
- remove .uuid
- [DONE] use context specific data structs
- [DONE] generalize context
- tests for handlers * cqrs & context
- [DONE] context with swapable model & read-write lock
- liveness probe / readiness probe
  - check kafka
  - remember last connection and let that be the test
  - the processor needs to signal, that it is stuck (msg unknown)
- /metrics endpoint for prometheus
  - add model size, writes, writetime, readtime, reads, handler calls to metrics
  - hook into sarama metrics to expose
- add timeouts & retries & exponential backoff & shedding & circuitbreaking & avoid thundering heard
- ensure ordered cart is exactly the cart shown in the browser
- better startup
  - snapshots
  - add a confirmation email
  - when triggering an sideeffect, then ensure to do this only once
  - after dirty shutdown wait for ip ttl
- use https://godoc.org/github.com/golang-collections/go-datastructures/slice/skip#SkipList.ByPosition
- debug
  - tracing via jaeger
  - debugger local
  - debug prod system
  - cpu & memory profiler
  - browser events in kafka
  - copy prod events to local

frontend

- production image
- remove .uuid
- use int32 everywhere https://package.elm-lang.org/packages/eriktim/elm-protocol-buffers/latest/#known-limitations
- use less & cleanup css classes from html elm
- [DONE] use protobuf / remove json
- [DONE] use modules
- use urls & links
- only process results matching the current model
- configure backend url
- tests
- add timeouts & retries & exponential backoff & shedding & circuitbreaking & avoid thundering heard

general

- create an uneasy environment https://github.com/Netflix/SimianArmy/tree/master/src/main/resources/scripts
- e2e tests
- add pim
- add fulfilment
- add cms
- add search

## Maybe later (stick to basics and prove of concept first)

- tilt

## Failed ideas

- use arc cache to skip marshal from struct to json
