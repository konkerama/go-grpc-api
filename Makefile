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

docker-deps:
	docker compose up

grpcui:
	grpcui -plaintext localhost:50051

grpcurl-list:
	grpcurl -plaintext localhost:50051 list

grpcurl-create-order:
	grpcurl -plaintext \
		-d '{"productName": "apples", "quantity": 5}' \
		localhost:50051 \
		helloworld.Orders.CreateOrder
