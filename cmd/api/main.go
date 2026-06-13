package main

import (
	"fmt"
	"hello-go/internal/handlers"
	"hello-go/internal/middleware"
	"hello-go/internal/models"
	"hello-go/internal/repository"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Conexão com Banco de Dados (SQLite para exemplo)
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("falha ao conectar no banco de dados")
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
