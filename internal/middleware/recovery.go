package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics and logs the error
func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					log.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("stack", string(debug.Stack())),
					)

					// Return 500 error
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"Internal Server Error","message":"An unexpected error occurred","code":500}`)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
