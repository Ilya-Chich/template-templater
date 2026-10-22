package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"rendering/internal/adapters/config"
)

const shutdownTimeout = 10 * time.Second

// Run serves until ctx is canceled, then shuts down within 10 seconds.
// Debug serves HTTP. Otherwise it serves TLS.
func Run(ctx context.Context, srv *http.Server, cfg config.Config) error {
	return serve(ctx, srv, func() error {
		if cfg.Debug {
			return srv.ListenAndServe()
		}
		return srv.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath)
	})
}

func serve(ctx context.Context, srv *http.Server, listen func() error) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- listen()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			return err
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
