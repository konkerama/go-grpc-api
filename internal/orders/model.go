package orders

import "time"

type Order struct {
	Id          string `json:"id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
	Status      string `json:"status"`
}

// OrderCreatedEvent represents the domain event payload
// that will be serialized (e.g., to JSON) and published to Kafka.
type OrderCreatedEvent struct {
	// Metadata for event tracking and idempotency
	EventID   string    `json:"event_id"`   // Unique UUID for this specific event instance
	EventType string    `json:"event_type"` // e.g., "orders.v1.order_created"
	Timestamp time.Time `json:"timestamp"`  // When the event occurred

	// Payload data reflecting the new state
	OrderID     string `json:"order_id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
	Status      string `json:"status"`
}
