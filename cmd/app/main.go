package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"clofi/internal/config"
	"clofi/internal/handler"
	"clofi/internal/repository/mongo"
	"clofi/internal/repository/postgres"
	"clofi/internal/service"
	authmw "clofi/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func serveFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	pgConnStr := "postgres://" + cfg.PostgresUser + ":" + cfg.PostgresPassword +
		"@" + cfg.PostgresHost + ":" + strconv.Itoa(cfg.PostgresPort) +
		"/" + cfg.PostgresDB + "?sslmode=disable"

	pgPool, err := pgxpool.New(context.Background(), pgConnStr)
	if err != nil {
		log.Fatalf("ошибка подключения к PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	mongoClient, err := mongodriver.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("ошибка подключения к MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(cfg.MongoDB)

	userRepo := postgres.NewUserRepository(pgPool)
	likeRepo := postgres.NewLikeRepository(pgPool)
	productRepo := mongo.NewProductRepository(mongoDB)

	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	likeService := service.NewLikeService(productRepo, userRepo, likeRepo)

	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret, cfg.JWTExpiresIn)
	productHandler := handler.NewProductHandler(productService)
	likeHandler := handler.NewLikeHandler(likeService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/products", productHandler.GetProducts)
		r.Get("/product", productHandler.GetProductByID)

		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthMiddleware(cfg.JWTSecret))
			r.Post("/like", likeHandler.ToggleLike)
		})
	})

	r.Get("/register", serveFile("./static/register/register.html"))
	r.Get("/login", serveFile("./static/login/login.html"))
	r.Get("/main", serveFile("./static/main/main.html"))
	r.Get("/favorites", serveFile("./static/favorites/favorites.html"))
	r.Get("/results", serveFile("./static/results/results.html"))

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Страница не найдена", http.StatusNotFound)
	})

	addr := cfg.ServerPort
	log.Printf("сервер запущен на http://localhost:%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
