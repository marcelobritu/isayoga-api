package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/database"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/repository/mongodb"
	"github.com/marcelobritu/isayoga-api/pkg/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	db, err := database.NewMongoDB(cfg.Database.MongoURI, cfg.Database.MongoDBName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer db.Client.Disconnect(context.Background())

	userRepo := mongodb.NewUserRepository(db.Database)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingAdmin, _ := userRepo.FindByEmail(ctx, "admin@isayoga.com")
	if existingAdmin != nil {
		log.Println("✅ Usuário admin já existe")
		log.Printf("   Email: admin@isayoga.com")
		log.Printf("   Senha: admin123")
		return
	}

	adminUser, err := entity.NewUser("Administrador", "admin@isayoga.com", "admin123", entity.RoleAdmin)
	if err != nil {
		log.Fatalf("Erro ao criar usuário admin: %v", err)
	}

	if err := userRepo.Create(ctx, adminUser); err != nil {
		log.Fatalf("Erro ao inserir usuário admin: %v", err)
	}

	log.Println("✅ Seed executado com sucesso!")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📋 CREDENCIAIS PADRÃO:")
	log.Println("   Email: admin@isayoga.com")
	log.Println("   Senha: admin123")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("⚠️  IMPORTANTE: Altere a senha após o primeiro login!")
	
	instructorUser, err := entity.NewUser("Instrutor Exemplo", "instrutor@isayoga.com", "instrutor123", entity.RoleInstructor)
	if err != nil {
		log.Printf("Aviso: Erro ao criar instrutor exemplo: %v", err)
		return
	}

	if err := userRepo.Create(ctx, instructorUser); err != nil {
		log.Printf("Aviso: Erro ao inserir instrutor exemplo: %v", err)
		return
	}

	log.Println("")
	log.Println("📋 INSTRUTOR EXEMPLO:")
	log.Println("   Email: instrutor@isayoga.com")
	log.Println("   Senha: instrutor123")
	
	studentUser, err := entity.NewUser("Aluno Exemplo", "aluno@isayoga.com", "aluno123", entity.RoleStudent)
	if err != nil {
		log.Printf("Aviso: Erro ao criar aluno exemplo: %v", err)
		return
	}

	if err := userRepo.Create(ctx, studentUser); err != nil {
		log.Printf("Aviso: Erro ao inserir aluno exemplo: %v", err)
		return
	}

	log.Println("")
	log.Println("📋 ALUNO EXEMPLO:")
	log.Println("   Email: aluno@isayoga.com")
	log.Println("   Senha: aluno123")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	fmt.Println("\n🎉 Dados iniciais criados com sucesso!")
}

