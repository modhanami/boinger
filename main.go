package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/hashers"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	log2 "log"
	"os"
)

func main() {
	parseFlagsAndEnvVars()
	db := initDB()
	gin.SetMode(gin.DebugMode)

	baseLogger := logger.NewZapLogger()
	boingServiceLogger := baseLogger.With("service", "boing")
	commentHandlerLogger := baseLogger.With("handler", "comment")

	userService := services.NewUserService(db, baseLogger)
	boingService := services.NewBoingService(db, boingServiceLogger)
	userTokenService := services.NewUserTokenService(db)
	authService := services.NewAuthService(db, userService, userTokenService, hashers.NewBcryptHasher())
	commentService := services.NewCommentService(db)
	//timelineService := services.NewTimelineService(userService, boingService)

	router := gin.Default()
	userTokenMiddleware := middlewares.MakeVerifyJWTUserTokenMiddleware(userTokenService)

	authGroup := router.Group("/auth")
	authGroup.POST("/register", endpoints.MakeRegisterEndpoint(authService))
	authGroup.POST("/login", endpoints.MakeLoginEndpoint(authService, userTokenService))
	authGroup.GET("/userinfo", userTokenMiddleware, endpoints.MakeUserInfoEndpoint())

	router.GET("/boings", endpoints.MakeListEndpoint(boingService))
	router.GET("/boings/:id", endpoints.MakeGetByIdEndpoint(boingService))
	router.POST("/boings", userTokenMiddleware, endpoints.MakeCreateEndpoint(boingService, userService))
	commentHandler := endpoints.NewCommentHandler(commentService, commentHandlerLogger)
	router.POST("/boings/:id/comments", userTokenMiddleware, commentHandler.Create)
	//router.GET("/timeline", endpoints.MakeTimelineEndpoint(timelineService))

	port := getEnv("PORT", "30027")
	address := "localhost:" + port
	router.Run(address)
}

func parseFlagsAndEnvVars() {
	flag.BoolVar(&endpoints.IsSecureCookieDisabled, "disable-secure-cookie", os.Getenv("DISABLE_SECURE_COOKIE") == "true", "disable secure cookie")
	flag.Parse()

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		log2.Fatalf("JWT_SECRET is not set")
	}
	services.JWTSecret = []byte(jwtSecret)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("dev.db?_foreign_keys=true"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Boing{}, &models.User{}, &models.RefreshToken{}, &models.Comment{})
	if err != nil {
		panic(err)
	}

	return db
}
