
codegen: frontend/src/Catalog.elm frontend/src/Checkout.elm backend/cmd/dev-tools/simba/simba
	find . -not -path ./backend/vendor -name '*.go' | xargs -L 1 go generate

backend/cmd/dev-tools/simba/simba: backend/cmd/dev-tools/simba/context.go.tpl backend/cmd/dev-tools/simba/main.go
	cd backend/cmd/dev-tools/simba && go generate && go build

frontend/src/Catalog.elm: backend/pkg/catalog/catalog.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/catalog/ catalog.proto

frontend/src/Checkout.elm: backend/pkg/checkout/checkout.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ --proto_path=backend/pkg/catalog/ checkout.proto
