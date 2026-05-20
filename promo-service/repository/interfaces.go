package repository

import "context"

type PromoRepository interface {
	CheckPromo(ctx context.Context, code string) (float64, error)
}