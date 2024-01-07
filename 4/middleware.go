package chiserver

import (
	"context"
	"net/http"
	"time"
)

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем таймаут в 1 секунду
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		// Передаем обновленный контекст в следующий обработчик
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
