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
	"net/http"
	"os"
)

func main() {
	parseFlagsAndEnvVars()
	db := initDB()
	gin.SetMode(gin.DebugMode)

	baseLogger := logger.NewZapLogger()
	boingServiceLogger := baseLogger.With("service", "boing")

	userService := services.NewUserService(db, baseLogger)
	boingService := services.NewBoingService(db, boingServiceLogger)
	userTokenService := services.NewUserTokenService(db)
	authService := services.NewAuthService(db, userService, userTokenService, hashers.NewBcryptHasher())
	//timelineService := services.NewTimelineService(userService, boingService)

	router := gin.Default()
	userTokenMiddleware := middlewares.MakeVerifyJWTUserTokenMiddleware(userTokenService)

	authGroup := router.Group("/auth")
	authGroup.POST("/register", endpoints.MakeRegisterEndpoint(authService))
	authGroup.POST("/login", endpoints.MakeLoginEndpoint(authService, userTokenService))
	authGroup.GET("/user-info", userTokenMiddleware, endpoints.MakeUserInfoEndpoint())

	router.GET("/boings", endpoints.MakeListEndpoint(boingService))
	router.GET("/boings/:id", endpoints.MakeGetByIdEndpoint(boingService))
	router.POST("/boings", userTokenMiddleware, endpoints.MakeCreateEndpoint(boingService, userService))
	//router.GET("/timeline", endpoints.MakeTimelineEndpoint(timelineService))

	router.POST("/dont-mind-me-boinging-around", userTokenMiddleware, func(c *gin.Context) {
		rawUserClaims, exists := c.Get(endpoints.UserClaimsKey)
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userClaims, ok := rawUserClaims.(*services.UserClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusOK, userClaims)
	})

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
	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
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
