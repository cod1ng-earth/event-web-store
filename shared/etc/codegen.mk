
codegen: frontend/src/Catalog.elm frontend/src/Checkout.elm backend/cmd/dev-tools/simba/simba
	find . -name '*.go' | grep -v vendor | grep -v shared | grep -v frontend | xargs -L 1 go generate
##	cd backend && ./cmd/dev-tools/simba/simba

backend/cmd/dev-tools/simba/simba: $(shell find backend/pkg/simba) backend/cmd/dev-tools/simba/main.go
	cd backend/pkg/simba && go generate
	cd backend/cmd/dev-tools/simba && go build

frontend/src/Catalog.elm: backend/pkg/context/catalog/api.proto
	true
##	protoc --elm_out=frontend/src --proto_path=backend/pkg/context/catalog/ api.proto

frontend/src/Checkout.elm: backend/pkg/context/checkout/api.proto
	true
##	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ api.proto


#protoc -I/usr/local/include -I. \
#  -I$GOPATH/src \
#  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#  --go_out=plugins=grpc:. \
#  api.proto
#
#protoc -I/usr/local/include -I. \
#  -I$GOPATH/src \
#  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#  --grpc-gateway_out=logtostderr=true:. \
#  api.proto
#
#protoc -I/usr/local/include -I. \
#  -I$GOPATH/src \
#  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#  --swagger_out=logtostderr=true:. \
#  api.proto
#
#npx openapi-generator generate -i ../../event-web-store/backend/pkg/contexts/catalog/api.swagger.json -g elm -o . #--additional-properties elmEnableCustomBasePaths=true
