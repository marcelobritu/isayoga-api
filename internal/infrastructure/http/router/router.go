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
	authHandler *handler.AuthHandler,
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
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
		})

		r.Route("/classes", func(r chi.Router) {
			r.Get("/", classHandler.List)
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AuthMiddleware)
				r.Use(customMiddleware.AdminOnly)
				r.Post("/", classHandler.Create)
			})
		})

		r.Route("/enrollments", func(r chi.Router) {
			r.Use(customMiddleware.AuthMiddleware)
			r.Post("/", enrollmentHandler.Enroll)
			r.Delete("/{id}", enrollmentHandler.Cancel)
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AuthMiddleware)
				r.Post("/change-password", userHandler.ChangePassword)
			})
			
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AuthMiddleware)
				r.Use(customMiddleware.AdminOnly)
				r.Get("/", userHandler.List)
				r.Post("/", userHandler.Create)
				r.Get("/{id}", userHandler.Get)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})
		})
	})

	return r
}
