package ports

import (
	"context"

	"rendering/internal/domain"
)

// Pages is the inbound port for page selection.
type Pages interface {
	Index(ctx context.Context) (domain.Page, error)
	Status(ctx context.Context, code int) (domain.Page, error)
}
