package hello

import "context"

type GreeterService struct {
}

func (s *GreeterService) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "Hello " + in.Name,
	}, nil
}
