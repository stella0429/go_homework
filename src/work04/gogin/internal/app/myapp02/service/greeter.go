package service

import (
	"context"

	v1 "gogin/api/myapp02/v1"
	"gogin/internal/app/myapp02/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Message: "Hello "}, nil
}
