package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_middleware "github.com/glebateee/todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	logger *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	cfg Config,
	logger *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) HTTPServer {
	return HTTPServer{
		mux:        http.NewServeMux(),
		config:     cfg,
		logger:     logger,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		s.logger.Warn("start HTTP server", zap.String("addr", s.config.Addr))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.logger.Warn("shutdown HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.logger.Warn("HTTP server stopped")
	}
	return nil
}
