package ports

import "context"

// Renderer executes a named template. It is the outbound port for HTML.
type Renderer interface {
	Render(ctx context.Context, name string, data any) ([]byte, error)
}
