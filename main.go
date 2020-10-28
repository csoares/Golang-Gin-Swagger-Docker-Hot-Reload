package main

import (
	"projetoapi/model"
	"projetoapi/routes"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var identityKey = "id"

func init() {
	services.OpenDatabase()
	services.Db.AutoMigrate(&model.Evaluation{})
	services.Db.AutoMigrate(&model.Users{})

	defer services.Db.Close()
}

func main() {

	services.FormatSwagger()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// NO AUTH
	router.GET("/echo/:echo", routes.EchoRepeat)

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
