
codegen: backend/pkg/catalog/catalog.pb.go backend/pkg/checkout/checkout.pb.go frontend/src/Catalog.elm frontend/src/Checkout.elm backend/pkg/catalog/context.go backend/pkg/checkout/context.go

backend/pkg/catalog/catalog.pb.go: backend/pkg/catalog/catalog.proto
	cd backend/pkg/catalog/ && protoc --go_out=. catalog.proto

backend/pkg/checkout/checkout.pb.go: backend/pkg/checkout/checkout.proto backend/pkg/catalog/catalog.proto
	cd backend/pkg/checkout/ && protoc --go_out=. --proto_path=../../../backend/pkg/checkout/ --proto_path=../../../backend/pkg/catalog/ checkout.proto

frontend/src/Catalog.elm: backend/pkg/catalog/catalog.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/catalog/ catalog.proto

frontend/src/Checkout.elm: backend/pkg/checkout/checkout.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ --proto_path=backend/pkg/catalog/ checkout.proto

backend/pkg/catalog/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/waitfree/context.go backend/pkg/catalog/catalog.proto
	cd backend && go run cmd/dev-tools/simba/main.go -readLock=wait-free -batch=true -name catalog > pkg/catalog/context.go
	gofmt -s -w backend/pkg/catalog/context.go

backend/pkg/checkout/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/waitfree/context.go backend/pkg/checkout/checkout.proto
	cd backend && go run cmd/dev-tools/simba/main.go -batch=false -name checkout > pkg/checkout/context.go
	gofmt -s -w backend/pkg/checkout/context.go
