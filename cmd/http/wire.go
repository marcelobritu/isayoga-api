//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/database"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/http/router"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/payment"
	mongoRepo "github.com/marcelobritu/isayoga-api/internal/infrastructure/repository/mongodb"
	"github.com/marcelobritu/isayoga-api/internal/interface/http/handler"
	"github.com/marcelobritu/isayoga-api/internal/usecase/class"
	enrollmentUC "github.com/marcelobritu/isayoga-api/internal/usecase/enrollment"
	paymentUC "github.com/marcelobritu/isayoga-api/internal/usecase/payment"
	"github.com/marcelobritu/isayoga-api/internal/usecase/user"
	"github.com/marcelobritu/isayoga-api/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeServer() (*Server, error) {
	wire.Build(
		config.Load,
		provideMongoDB,
		provideMongoDatabase,
		provideMongoClient,
		provideUserRepository,
		provideClassRepository,
		provideEnrollmentRepository,
		providePaymentRepository,
		provideMercadoPagoClient,
		user.NewCreateUserUseCase,
		user.NewGetUserUseCase,
		user.NewListUsersUseCase,
		user.NewUpdateUserUseCase,
		user.NewDeleteUserUseCase,
		class.NewCreateClassUseCase,
		class.NewListClassesUseCase,
		enrollmentUC.NewEnrollStudentUseCase,
		enrollmentUC.NewCancelEnrollmentUseCase,
		paymentUC.NewProcessWebhookUseCase,
		handler.NewHealthHandler,
		handler.NewUserHandler,
		handler.NewClassHandler,
		handler.NewEnrollmentHandler,
		handler.NewWebhookHandler,
		router.Setup,
		NewServer,
	)
	return &Server{}, nil
}

func provideMongoDB(cfg *config.Config) (*database.MongoDB, error) {
	return database.NewMongoDB(cfg.Database.MongoURI, cfg.Database.MongoDBName)
}

func provideMongoDatabase(mongodb *database.MongoDB) *mongo.Database {
	return mongodb.Database
}

func provideMongoClient(mongodb *database.MongoDB) *mongo.Client {
	return mongodb.Client
}

func provideUserRepository(db *mongo.Database) repository.UserRepository {
	return mongoRepo.NewUserRepository(db)
}

func provideClassRepository(db *mongo.Database, client *mongo.Client) repository.ClassRepository {
	return mongoRepo.NewClassRepository(db, client)
}

func provideEnrollmentRepository(db *mongo.Database) repository.EnrollmentRepository {
	return mongoRepo.NewEnrollmentRepository(db)
}

func providePaymentRepository(db *mongo.Database) repository.PaymentRepository {
	return mongoRepo.NewPaymentRepository(db)
}

func provideMercadoPagoClient(cfg *config.Config) *payment.MercadoPagoClient {
	return payment.NewMercadoPagoClient(cfg.MercadoPago.AccessToken)
}
