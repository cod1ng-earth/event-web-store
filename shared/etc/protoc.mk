

protoc: backend/pkg/catalog/catalog.pb.go backend/pkg/checkout/checkout.pb.go frontend/src/Catalog.elm frontend/src/Checkout.elm

backend/pkg/catalog/catalog.pb.go: backend/pkg/catalog/catalog.proto
	cd backend/pkg/catalog/ && protoc --go_out=. catalog.proto

backend/pkg/checkout/checkout.pb.go: backend/pkg/checkout/checkout.proto backend/pkg/catalog/catalog.proto
	cd backend/pkg/checkout/ && protoc --go_out=. --proto_path=../../../backend/pkg/checkout/ --proto_path=../../../backend/pkg/catalog/ checkout.proto

frontend/src/Catalog.elm: backend/pkg/catalog/catalog.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/catalog/ catalog.proto

frontend/src/Checkout.elm: backend/pkg/checkout/checkout.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ --proto_path=backend/pkg/catalog/ checkout.proto
