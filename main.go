package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"sci-review/auth"
	"sci-review/user"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("starting application")

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("environment variables loaded")

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	slog.Info("set application environment", "appEnv", appEnv)

	dataSourceName := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("database connected established")

	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	loginAttemptRepo := auth.NewLoginAttemptRepo(db)
	authService := auth.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)
	slog.Info("services initialized")

	r := gin.Default()
	auth.Register(r, authService)
	user.Register(r, userService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/protected", auth.JwtMiddleware(), func(c *gin.Context) {
		principal := c.MustGet("principal").(*auth.Principal)
		c.JSON(http.StatusOK, gin.H{
			"message":   "protected route",
			"principal": principal,
		})
	})
	slog.Info("routes registered")

	r.Run(os.Getenv("PORT"))
}
