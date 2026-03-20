package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gorm.io/driver/postgres"

	"projetoapi/docs"
	"projetoapi/internal/domain"
	"projetoapi/internal/handler"
	"projetoapi/internal/infrastructure/database"
	"projetoapi/internal/infrastructure/jwt"
	"projetoapi/internal/middleware"
	"projetoapi/internal/repository"
	"projetoapi/internal/service"
)

// @title API de avaliações
// @version 1.0
// @description Essa API permite manter todas as avaliações realizadas.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize database
	dbConfig, err := database.ReadConfig("config/db.config")
	if err != nil {
		log.Fatalf("Failed to read database config: %v", err)
	}

	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	db.AutoMigrate(&domain.Evaluation{})
	db.AutoMigrate(&domain.User{})

	// Create audit_logs table for transaction demo
	db.Exec(`CREATE TABLE IF NOT EXISTS audit_logs (
		id SERIAL PRIMARY KEY,
		action VARCHAR(100),
		resource_id INTEGER,
		created_at TIMESTAMP
	)`)

	// Seed database
	seedDatabase(db)

	// Initialize services
	jwtService, err := jwt.NewService("config/secretKey.key")
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	// Initialize repositories
	evaluationRepo := repository.NewEvaluationRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Initialize services
	evaluationService := service.NewEvaluationService(evaluationRepo)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	evaluationHandler := handler.NewEvaluationHandler(evaluationService)
	authHandler := handler.NewAuthHandler(userService, jwtService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Configure Swagger
	docs.SwaggerInfo.Title = "API de avaliações"
	docs.SwaggerInfo.Description = "Essa API permite manter todas as avaliações realizadas."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Create router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// NO AUTH
	router.GET("/api/v1/echo", evaluationHandler.Echo)

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// Evaluation routes (AUTH required)
	evaluation := router.Group("/api/v1/evaluation")
	evaluation.Use(authMiddleware)
	{
		evaluation.POST("/", evaluationHandler.Create)
		evaluation.GET("/", evaluationHandler.GetAll)
		evaluation.GET("/raw", evaluationHandler.GetByRawQuery)
		evaluation.PUT("/batch", evaluationHandler.UpdateBatch)
		evaluation.POST("/audit", evaluationHandler.CreateWithAudit)
		evaluation.GET("/:id", evaluationHandler.GetByID)
		evaluation.PUT("/:id", evaluationHandler.Update)
		evaluation.DELETE("/:id", evaluationHandler.Delete)
	}

	// Auth routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.PUT("/refresh_token", authMiddleware, authHandler.RefreshToken)
	}

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	router.Run(":8080")
}

func seedDatabase(db *database.DB) {
	// Seed default user
	var count int64
	db.Model(&domain.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash := service.HashPasswordForSeeding("admin123")
		db.Create(&domain.User{Username: "admin", Password: hash})
	}

	// Seed default evaluations
	db.Model(&domain.Evaluation{}).Count(&count)
	if count == 0 {
		evaluations := []domain.Evaluation{
			{Rating: 5, Note: "Excellent service, very satisfied!"},
			{Rating: 4, Note: "Good experience overall."},
			{Rating: 3, Note: "Average, could be improved."},
			{Rating: 2, Note: "Below expectations."},
			{Rating: 1, Note: "Very disappointing experience."},
		}
		db.Create(&evaluations)
	}
}
