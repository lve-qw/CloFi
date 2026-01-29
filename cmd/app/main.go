// Пакет main — точка входа приложения.
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
	authmw "clofi/pkg/middleware" // ✅ Наш JWT middleware

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware" // ✅ chi middleware
	"github.com/jackc/pgx/v5/pgxpool"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	// PostgreSQL
	pgConnStr := "postgres://" + cfg.PostgresUser + ":" + cfg.PostgresPassword +
		"@" + cfg.PostgresHost + ":" + strconv.Itoa(cfg.PostgresPort) +
		"/" + cfg.PostgresDB + "?sslmode=disable"

	pgPool, err := pgxpool.New(context.Background(), pgConnStr)
	if err != nil {
		log.Fatalf("ошибка подключения к PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	// MongoDB
	mongoClient, err := mongodriver.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("ошибка подключения к MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(cfg.MongoDB)

	// Репозитории
	userRepo := postgres.NewUserRepository(pgPool)    // ✅ Правильное имя
	likeRepo := postgres.NewLikeRepository(pgPool)    // ✅
	productRepo := mongo.NewProductRepository(mongoDB) // ✅

	// Сервисы
	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	likeService := service.NewLikeService(productRepo, userRepo, likeRepo)

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret, cfg.JWTExpiresIn)
	productHandler := handler.NewProductHandler(productService)
	likeHandler := handler.NewLikeHandler(likeService)

	// Роутер
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // из chi
	r.Use(middleware.Recoverer) // из chi

	// Публичные маршруты
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/product", productHandler.GetProductByID)

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(authmw.AuthMiddleware(cfg.JWTSecret)) // ✅ Используем алиас authmw
		r.Post("/like", likeHandler.ToggleLike)
	})

	addr := ":" + cfg.ServerPort
	log.Printf("сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}


