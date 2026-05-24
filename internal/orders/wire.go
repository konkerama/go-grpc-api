package orders

import "github.com/konkerama/go-grpc-api/internal/db"

type Module struct {
	Controller *Controller
	Service    *Service
}

func Wire(pool db.DBPool) *Module {
	repo := NewRepository(pool)
	svc := NewService(repo)
	ctl := NewController(svc)
	return &Module{Controller: ctl, Service: svc}
}
