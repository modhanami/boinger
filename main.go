package main

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
)

func main() {
	db := initDB()

	boingService := services.NewBoingService(db)
	userService := services.NewUserService(db)
	userTokenService := services.NewUserTokenService()
	authService := services.NewAuthService(db, userService, userTokenService)

	router := gin.Default()
	userTokenMiddleware := middlewares.MakeVerifyJWTUserTokenMiddleware(userTokenService)

	authGroup := router.Group("/auth")
	authGroup.POST("/register", endpoints.MakeRegisterEndpoint(authService))
	authGroup.POST("/login", endpoints.MakeLoginEndpoint(authService))

	router.GET("/boings", endpoints.MakeListEndpoint(boingService))
	router.GET("/boings/:id", endpoints.MakeGetByIdEndpoint(boingService))
	router.POST("/boings", endpoints.MakeCreateEndpoint(boingService))

	router.POST("/dont-mind-me-boinging-around", userTokenMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userId": c.GetString(middlewares.UserIdKey),
		})
	})

	port := getEnv("PORT", "30027")
	address := "localhost:" + port
	router.Run(address)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.BoingModel{}, &models.UserModel{})
	if err != nil {
		panic(err)
	}

	return db
}
