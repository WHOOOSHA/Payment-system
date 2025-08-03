package router

import (
	"net/http"
	"payment_system/internal/handlers"

	"payment_system/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(secret string, useAuth bool) http.Handler {
	r := chi.NewRouter()

	//Проверка подписи совершается только в методе,
	//в котором изменяются значения в бд
	if useAuth {
		r.With(middleware.HMACAuthMiddleware(secret)).Post("/api/send", handlers.Send)
	} else {
		r.Post("/api/send", handlers.Send)
	}

	r.Get("/api/transactions", handlers.GetLastTransactions)
	r.Get("/api/wallet/{address}/balance", handlers.GetBalance)

	return r
}
