package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIDHeader)
			if requestId == "" {
				requestId = uuid.NewString()
				r.Header.Set(requestIDHeader, requestId)
			}
			w.Header().Set(requestIDHeader, requestId)
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIDHeader)
			l := logger.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if arg := recover(); arg != nil {
					logger := core_logger.FromContextMust(r.Context())
					response_handler := core_http_response.NewHTTPResponseHandler(logger, w)
					response_handler.PanicResponse(
						arg,
						"unexpected panic during handle HTTP request",
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContextMust(ctx)
			codeWriter := core_http_response.NewWriter(w)
			start := time.Now()
			logger.Debug(">>> incoming HTTP request", zap.String("http_method", r.Method), zap.Time("time", start.UTC()))

			next.ServeHTTP(codeWriter, r)

			logger.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", codeWriter.StatusCode()),
				zap.Duration("latency", time.Since(start)),
			)

		})
	}
}
