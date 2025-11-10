package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	pkgAuth "github.com/marcelobritu/isayoga-api/pkg/auth"
	"github.com/marcelobritu/isayoga-api/pkg/config"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"github.com/marcelobritu/isayoga-api/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Router *chi.Mux
}

func NewServer(cfg *config.Config, r *chi.Mux) *Server {
	return &Server{
		Config: cfg,
		Router: r,
	}
}

func main() {
	srv, err := InitializeServer()
	if err != nil {
		fmt.Printf("Erro ao inicializar servidor: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(srv.Config.Server.Env); err != nil {
		fmt.Printf("Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	pkgAuth.SetJWTSecret(srv.Config.Auth.JWTSecret)

	tp, shutdown, err := telemetry.InitTracer(telemetry.Config{
		ServiceName:    srv.Config.Telemetry.ServiceName,
		ServiceVersion: srv.Config.Telemetry.ServiceVersion,
		ZipkinURL:      srv.Config.Telemetry.ZipkinURL,
		Environment:    srv.Config.Server.Env,
	})
	if err != nil {
		logger.Warn("Falha ao inicializar tracing", zap.Error(err))
	} else {
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				logger.Error("Erro ao finalizar tracing", zap.Error(err))
			}
		}()
		logger.Info("OpenTelemetry inicializado com sucesso",
			zap.String("exporter", "zipkin"),
			zap.String("service", srv.Config.Telemetry.ServiceName),
		)
	}

	_ = tp

	logger.Info("Iniciando API IsaYoga",
		zap.String("environment", srv.Config.Server.Env),
		zap.String("version", srv.Config.Telemetry.ServiceVersion),
		zap.String("architecture", "Clean Architecture"),
	)

	// Debug: Print all registered routes
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logger.Info("Rotas Registradas:")
	chi.Walk(srv.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*", "", -1)
		logger.Info("Rota", zap.String("method", method), zap.String("path", route))
		return nil
	})
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	addr := fmt.Sprintf("%s:%s", srv.Config.Server.Host, srv.Config.Server.Port)

	logger.Info("Servidor iniciado com sucesso",
		zap.String("address", addr),
		zap.String("health_check", fmt.Sprintf("http://%s/health", addr)),
		zap.String("api_endpoint", fmt.Sprintf("http://%s/api/v1", addr)),
		zap.String("zipkin_ui", "http://localhost:9411"),
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	handler := otelhttp.NewHandler(srv.Router, "isayoga-api",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
	)

	go func() {
		if err := http.ListenAndServe(addr, handler); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Encerrando servidor graciosamente...")
}
