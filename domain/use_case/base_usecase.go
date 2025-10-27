package use_case

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type UseCase[RequestData any, ResponseData any] interface {
	Validate(fiberCtx *fiber.Ctx, ctx context.Context, input RequestData) error
	Execute(fiberCtx *fiber.Ctx, ctx context.Context, input RequestData) (ResponseData, error)
	Transform(fiberCtx *fiber.Ctx, ctx context.Context, result ResponseData) (any, error)
}
