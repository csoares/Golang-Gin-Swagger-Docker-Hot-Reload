package main

import (
	"projetoapi/controllers"
	"projetoapi/model"
	"projetoapi/routes"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/postgres"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var identityKey = "id"

func init() {
	services.OpenDatabase()
	services.Db.AutoMigrate(&model.Evaluation{})
	services.Db.AutoMigrate(&model.User{})
	seedDatabase()
}

func seedDatabase() {
	// Seed default user
	var count int64
	services.Db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash, _ := controllers.HashPassword("admin123")
		services.Db.Create(&model.User{Username: "admin", Password: hash})
	}

	// Seed default evaluations
	services.Db.Model(&model.Evaluation{}).Count(&count)
	if count == 0 {
		evaluations := []model.Evaluation{
			{Rating: 5, Note: "Excellent service, very satisfied!"},
			{Rating: 4, Note: "Good experience overall."},
			{Rating: 3, Note: "Average, could be improved."},
			{Rating: 2, Note: "Below expectations."},
			{Rating: 1, Note: "Very disappointing experience."},
		}
		services.Db.Create(&evaluations)
	}
}

func main() {

	services.FormatSwagger()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// NO AUTH
	router.GET("/api/v1/echo", routes.EchoRepeat)

	// AUTH
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	evaluation := router.Group("/api/v1/evaluation")
	evaluation.Use(services.AuthorizationRequired())
	{
		evaluation.POST("/", routes.AddEvaluation)
		evaluation.GET("/", routes.GetAllEvaluation)
		evaluation.GET("/:id", routes.GetEvaluationById)
		evaluation.PUT("/:id", routes.UpdateEvaluation)
		evaluation.DELETE("/:id", routes.DeleteEvaluation)
	}

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", routes.GenerateToken)
		auth.POST("/register", routes.RegisterUser)
		auth.PUT("/refresh_token", services.AuthorizationRequired(), routes.RefreshToken)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
