package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/marcelobritu/isayoga-api/internal/interface/http/handler"
	customMiddleware "github.com/marcelobritu/isayoga-api/internal/infrastructure/http/middleware"
)

func Setup(
	healthHandler *handler.HealthHandler,
	userHandler *handler.UserHandler,
	classHandler *handler.ClassHandler,
	enrollmentHandler *handler.EnrollmentHandler,
	webhookHandler *handler.WebhookHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", healthHandler.Check)

	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/mercadopago", webhookHandler.MercadoPago)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)
			r.Get("/{id}", userHandler.Get)
			r.Put("/{id}", userHandler.Update)
			r.Delete("/{id}", userHandler.Delete)
		})

		r.Route("/classes", func(r chi.Router) {
			r.Post("/", classHandler.Create)
			r.Get("/", classHandler.List)
		})

		r.Route("/enrollments", func(r chi.Router) {
			r.Post("/", enrollmentHandler.Enroll)
			r.Delete("/{id}", enrollmentHandler.Cancel)
		})
	})

	return r
}
