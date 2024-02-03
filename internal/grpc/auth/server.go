package auth

import (
	"github.com/DemianShtepa/exception-control/protos/gen"
	"google.golang.org/grpc"
)

type serverApi struct {
	gen.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	gen.RegisterAuthServer(gRPC, &serverApi{})
}
