package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaw3-q/api-blog-go/internal/handlers"
	"github.com/kaw3-q/api-blog-go/internal/middleware"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, lendo variáveis de ambiente do sistema")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL não configurada no ambiente")
	}

	// Conexão com Banco de Dados PostgreSQL não local
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("falha ao conectar no banco de dados: " + err.Error())
	}

	// Auto-Migrate (Cria tabelas automaticamente)
	db.AutoMigrate(&models.User{}, &models.Post{})

	// Repositórios SQL (Proteção contra SQL Injection nativa via GORM/Prepared Statements)
	postRepo := repository.NewSQLPostRepository(db)
	userRepo := repository.NewSQLUserRepository(db)

	// Handlers
	postHandler := handlers.NewPostHandler(postRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

	mux := http.NewServeMux()

	// Rotas Públicas
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/register", authHandler.Register)

	// Rotas Protegidas (Exemplo de uso de Middleware)
	// Para manter o exemplo simples com http.NewServeMux, aplicamos o middleware manualmente
	// Em frameworks como Gin ou Echo, isso seria mais elegante.
	
	mux.Handle("/posts", middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(postHandler.ServeHTTP)).ServeHTTP(w, r)
		} else {
			postHandler.ServeHTTP(w, r)
		}
	})))

	mux.Handle("/admin/users", middleware.AuthMiddleware(middleware.RoleMiddleware(models.RoleAdmin)(userHandler)))

	port := ":8080"
	fmt.Printf("Servidor Blog SEGURO rodando em http://localhost%s\n", port)
	http.ListenAndServe(port, mux)
}
