package biz

import (
	"context"
)

type Greeter struct {
	Hello string
}

type GreeterRepo interface {
	CreateGreeter(context.Context, *Greeter) error
	UpdateGreeter(context.Context, *Greeter) error
}

type GreeterUsecase struct {
	repo GreeterRepo
}

func NewGreeterUsecase(repo GreeterRepo) *GreeterUsecase {
	return &GreeterUsecase{repo: repo}
}

func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) error {
	return uc.repo.CreateGreeter(ctx, g)
}

func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) error {
	return uc.repo.UpdateGreeter(ctx, g)
}
