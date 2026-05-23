# go-grpc-api


## todo:
- add proper orders endpoints
- add users endpoints (see how they interact)
- add security (TLS and auth) (By default, gRPC expects production-grade TLS (SSL encryption) to protect the binary stream. Because you are testing this locally on your own machine, this flag explicitly tells the client: "It's okay, turn off TLS validation for this connection.")
- Add proper orders, suppor controller + service + repository in posgres
- POC and introduce Kafka





Kafka architecture:
```
[ Python Client ] or [ Go Client ]
        │
        │ (Synchronous gRPC Call: "CreateOrder")
        ▼
┌─────────────────────────────────┐
│          Order Service          │  ──► [ PostgreSQL ] (Save Order State)
│          (gRPC Server)          │
└─────────────────────────────────┘
        │
        │ (Asynchronous Event: "order-created")
        ▼
   █████████████████████████████ 
   █       KAFKA TOPIC         █   ◄── [ Event Broker ]
   █████████████████████████████ 
        │
        │ (Streaming Consumer)
        ▼
┌─────────────────────────────────┐
│     Rewards/Loyalty Service     │  ──► [ Processes tokens or logs points ]
│       (Kafka Consumer)          │
└─────────────────────────────────┘
```


project structure
```
grpc-kafka-poc/
├── api/
│   └── proto/
│       └── orders/
│           └── v1/
│               └── orders.proto       # Core data contract shared by everything
├── cmd/
│   ├── order-service/                 # Binary 1: Runs the gRPC Server + Kafka Producer
│   │   └── main.go
│   └── rewards-service/               # Binary 2: Runs the standalone Kafka Consumer
│       └── main.go
├── internal/
│   ├── kafka/                         # Helpers to initialize Kafka writers/readers
│   │   ├── producer.go
│   │   └── consumer.go
│   └── orders/
│       └── handler.go                 # Implements gRPC logic & triggers Kafka messages
├── pkg/
│   └── pb/                            # Auto-generated Go files from protoc
│       └── orders/
│           └── v1/
└── docker-compose.yaml                # Spins up Kafka & Zookeeper locally with 1 command
```