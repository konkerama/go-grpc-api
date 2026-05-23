.PHONY: proto

proto:
	protoc --proto_path=api/proto \
		--go_out=pkg/pb --go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
		api/proto/orders/v1/orders.proto

run-server:
	go run ./cmd/server/.

run-client:
	go run ./cmd/client/.
