package router

import (
	"github.com/gleb-korostelev/gophermart.git/internal/middleware"
	"github.com/gleb-korostelev/gophermart.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RouterInit(svc service.APIServiceI, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	router.Post("/api/user/register", svc.Register)
	router.Post("/api/user/login", svc.Login)
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.EnsureUserCookie)
		r.Post("/api/user/orders", svc.Orders)
		r.Post("/api/user/balance/withdraw", svc.Withdraw)
		r.Get("/api/user/orders", svc.GetOrders)
		r.Get("/api/user/balance", svc.GetBalance)
		r.Get("/api/user/withdrawals", svc.GetWithdrawals)
	})

	return router
}
