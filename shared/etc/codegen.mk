
codegen: backend/pkg/catalog/catalog.pb.go backend/pkg/checkout/checkout.pb.go frontend/src/Catalog.elm frontend/src/Checkout.elm backend/pkg/catalog/context.go backend/pkg/checkout/context.go backend/pkg/pim/context.go backend/pkg/warehouse/context.go backend/pkg/warehouse/warehouse.pb.go backend/pkg/pim/pim.pb.go

backend/pkg/catalog/catalog.pb.go: backend/pkg/catalog/catalog.proto
	cd backend/pkg/catalog/ && protoc --go_out=. catalog.proto

backend/pkg/checkout/checkout.pb.go: backend/pkg/checkout/checkout.proto
	cd backend/pkg/checkout/ && protoc --go_out=. checkout.proto

frontend/src/Catalog.elm: backend/pkg/catalog/catalog.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/catalog/ catalog.proto

frontend/src/Checkout.elm: backend/pkg/checkout/checkout.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ --proto_path=backend/pkg/catalog/ checkout.proto

backend/pkg/catalog/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/context.go.tpl backend/pkg/catalog/catalog.proto
	cd backend && go run cmd/dev-tools/simba/main.go -readLock=wait-free -batch=true -name catalog > pkg/catalog/context.go
	gofmt -s -w backend/pkg/catalog/context.go

backend/pkg/checkout/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/context.go.tpl backend/pkg/checkout/checkout.proto
	cd backend && go run cmd/dev-tools/simba/main.go -batch=false -name checkout > pkg/checkout/context.go
	gofmt -s -w backend/pkg/checkout/context.go

backend/pkg/pim/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/context.go.tpl backend/pkg/pim/pim.proto
	cd backend && go run cmd/dev-tools/simba/main.go -batch=false -name pim > pkg/pim/context.go
	gofmt -s -w backend/pkg/pim/context.go

backend/pkg/warehouse/context.go: backend/cmd/dev-tools/simba/main.go backend/cmd/dev-tools/simba/context.go.tpl backend/pkg/warehouse/warehouse.proto
	cd backend && go run cmd/dev-tools/simba/main.go -batch=false -name warehouse > pkg/warehouse/context.go
	gofmt -s -w backend/pkg/warehouse/context.go

backend/pkg/warehouse/warehouse.pb.go: backend/pkg/warehouse/warehouse.proto
	cd backend/pkg/warehouse/ && protoc --go_out=. warehouse.proto

backend/pkg/pim/pim.pb.go: backend/pkg/pim/pim.proto
	cd backend/pkg/pim/ && protoc --go_out=. pim.proto
