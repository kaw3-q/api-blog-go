package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kaw3-q/api-blog-go/internal/handlers"
	"github.com/kaw3-q/api-blog-go/internal/middleware"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"
	"github.com/kaw3-q/api-blog-go/prisma/db"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, lendo variáveis de ambiente do sistema")
	}

	// Inicialização do Prisma Client Go
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		panic("falha ao conectar no banco de dados via Prisma: " + err.Error())
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			log.Printf("erro ao desconectar o Prisma Client: %v", err)
		}
	}()

	// Repositórios usando Prisma Client
	userRepo := repository.NewPrismaUserRepository(client)
	postRepo := repository.NewPrismaPostRepository(client)
	categoryRepo := repository.NewPrismaCategoryRepository(client)
	tagRepo := repository.NewPrismaTagRepository(client)
	commentRepo := repository.NewPrismaCommentRepository(client)
	likeRepo := repository.NewPrismaLikeRepository(client)

	// Handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	tagHandler := handlers.NewTagHandler(tagRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)
	likeHandler := handlers.NewLikeHandler(likeRepo)

	mux := http.NewServeMux()

	// Rotas Públicas de Autenticação
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/register", authHandler.Register)

	// Rota de Posts (Pública para GET, Protegida via AuthMiddleware para POST)
	mux.Handle("/posts", middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(postHandler.ServeHTTP)).ServeHTTP(w, r)
		} else {
			postHandler.ServeHTTP(w, r)
		}
	})))

	// Rota de Categorias (Pública para GET, Protegida via AuthMiddleware para POST)
	mux.Handle("/categories", middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(categoryHandler.ServeHTTP)).ServeHTTP(w, r)
		} else {
			categoryHandler.ServeHTTP(w, r)
		}
	})))

	// Rota de Tags (Pública para GET, Protegida via AuthMiddleware para POST)
	mux.Handle("/tags", middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(tagHandler.ServeHTTP)).ServeHTTP(w, r)
		} else {
			tagHandler.ServeHTTP(w, r)
		}
	})))

	// Rota de Comentários (Pública para GET/Listagem, Protegida via AuthMiddleware para POST/Criação)
	mux.Handle("/comments", middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(commentHandler.ServeHTTP)).ServeHTTP(w, r)
		} else {
			commentHandler.ServeHTTP(w, r)
		}
	})))

	// Rota de Likes (Totalmente protegida - apenas usuários logados curtem)
	mux.Handle("/likes", middleware.SecurityHeadersMiddleware(middleware.AuthMiddleware(likeHandler)))

	// Rotas de Administração (Apenas Admin)
	mux.Handle("/admin/users", middleware.AuthMiddleware(middleware.RoleMiddleware(models.RoleAdmin)(userHandler)))

	port := ":8080"
	fmt.Printf("Servidor Blog com Prisma rodando em http://localhost%s\n", port)
	if err := http.ListenAndServe(port, middleware.CORSMiddleware(mux)); err != nil {
		log.Fatalf("erro ao iniciar o servidor: %v", err)
	}
}
