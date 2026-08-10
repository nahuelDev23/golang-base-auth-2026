package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"test.nahueldev23.com/internal/auth"
	"test.nahueldev23.com/internal/config"
	"test.nahueldev23.com/internal/session"
	"test.nahueldev23.com/internal/user"
)

func New(db *pgxpool.Pool, cfg *config.Config) *gin.Engine {

	router := gin.Default()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	sessionRepository := session.NewRepository(db)

	authService := auth.NewService(
		userRepository,
		sessionRepository,
		cfg.JWTSecret, // esto lo ajustaremos en un momento
	)

	authHandler := auth.NewHandler(authService)

	router.POST("/login", authHandler.Login)
	router.GET(
		"/users/:username",
		userHandler.GetByUsername,
	)

	//la app esta viva?
	router.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status": "ok",
		})

	})

	//la app puede recibir trafico?
	router.GET("/ready", func(c *gin.Context) {

		err := db.Ping(c.Request.Context())

		if err != nil {
			c.JSON(503, gin.H{
				"database": "unavailable",
			})
			return
		}

		c.JSON(200, gin.H{
			"database": "ready",
		})
	})

	router.POST(
		"/register",
		userHandler.Register,
	)

	router.GET(
		"/me",
		authService.AuthMiddleware(),
		userHandler.Me,
	)

	router.POST(
		"/logout",
		authService.AuthMiddleware(),
		authHandler.Logout,
	)

	router.POST(
		"/logout-all",
		authService.AuthMiddleware(),
		authHandler.LogoutAll,
	)
	return router
}
