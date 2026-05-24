package orders

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/konkerama/go-grpc-api/pkg/pb/orders/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersService interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
}

type Controller struct {
	service *Service
	pb.UnimplementedOrdersServer
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// CreateOrder handles the incoming gRPC request for order creation.
func (c *Controller) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	// Structured, context-aware info log
	slog.InfoContext(ctx, "received gRPC request to create order",
		slog.String("product_name", in.GetProductName()),
		slog.Int64("quantity", in.GetQuantity()),
	)

	// 1. Transport-Level Validation
	if in.GetProductName() == "" {
		slog.WarnContext(ctx, "validation failed: empty product name")
		return nil, status.Error(codes.InvalidArgument, "product name cannot be empty")
	}
	if in.GetQuantity() <= 0 {
		slog.WarnContext(ctx, "validation failed: invalid quantity", slog.Int64("quantity", in.GetQuantity()))
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	// 2. Data Mapping: Protobuf Request -> Domain Model
	domainOrder := &Order{
		ProductName: in.GetProductName(),
		Quantity:    in.GetQuantity(),
	}

	// 3. Invoke the Business Service Layer
	createdOrder, err := c.service.CreateOrder(ctx, domainOrder)
	if err != nil {
		return nil, c.mapError(ctx, err)
	}

	// 4. Data Mapping: Domain Model -> Protobuf Reply
	return &pb.CreateOrderReply{
		OrderID: createdOrder.Id,
	}, nil
}

// mapError translates internal business domain errors and logs them with context.
func (c *Controller) mapError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		slog.WarnContext(ctx, "client canceled the request mid-flight")
		return status.Error(codes.Canceled, "request was canceled")

	case errors.Is(err, context.DeadlineExceeded):
		slog.ErrorContext(ctx, "request timed out before completion", slog.Any("error", err))
		return status.Error(codes.DeadlineExceeded, "request timed out")

	default:
		// Log the structural details of the raw internal error, but mask it from the final client reply
		slog.ErrorContext(ctx, "internal service execution failure",
			slog.Any("error", err),
		)
		return status.Error(codes.Internal, "an internal error occurred")
	}
}
