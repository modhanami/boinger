package main

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/hashers"
	"github.com/modhanami/boinger/log"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
)

func main() {
	db := initDB()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	gin.SetMode(gin.DebugMode)

	zapSugaredLogger := zapLogger.Sugar()

	userService := services.NewUserService(db, log.NewZapLoggerAdapter(zapSugaredLogger.With("service", "user")))
	boingService := services.NewBoingService(db)
	userTokenService := services.NewUserTokenService(db)
	authService := services.NewAuthService(db, userService, userTokenService, hashers.NewBcryptHasher())
	timelineService := services.NewTimelineService(userService, boingService)

	router := gin.Default()
	userTokenMiddleware := middlewares.MakeVerifyJWTUserTokenMiddleware(userTokenService)

	authGroup := router.Group("/auth")
	authGroup.POST("/register", endpoints.MakeRegisterEndpoint(authService))
	authGroup.POST("/login", endpoints.MakeLoginEndpoint(authService, userTokenService))
	authGroup.GET("/user-info", userTokenMiddleware, endpoints.MakeUserInfoEndpoint())

	router.GET("/boings", endpoints.MakeListEndpoint(boingService))
	router.GET("/boings/:id", endpoints.MakeGetByIdEndpoint(boingService))
	router.POST("/boings", userTokenMiddleware, endpoints.MakeCreateEndpoint(boingService, userService))
	router.GET("/timeline", endpoints.MakeTimelineEndpoint(timelineService))

	router.POST("/dont-mind-me-boinging-around", userTokenMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userId": c.GetString(endpoints.UserIdKey),
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

	err = db.AutoMigrate(&models.Boing{}, &models.User{}, &models.RefreshToken{})
	if err != nil {
		panic(err)
	}

	return db
}
