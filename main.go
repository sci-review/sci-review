package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"os"
	"sci-review/handler"
	"sci-review/repo"
	"sci-review/service"
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

	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo)
	organizationRepo := repo.NewOrganizationRepo(db)
	organizationService := service.NewOrganizationService(organizationRepo)
	reviewRepo := repo.NewReviewRepo(db)
	reviewService := service.NewReviewService(reviewRepo)
	preliminaryInvestigationRepo := repo.NewPreliminaryInvestigationRepo(db)
	preliminaryInvestigationService := service.NewPreliminaryInvestigationService(preliminaryInvestigationRepo)
	slog.Info("services initialized")

	authMiddleware := handler.AuthMiddleware()
	slog.Info("middleware initialized")

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/assets", "./assets")

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))

	handler.RegisterHomeHandler(r, authMiddleware)
	handler.RegisterAuthHandler(r, authService)
	handler.RegisterUserHandler(r, userService)
	handler.RegisterOrganizationHandler(r, organizationService, authMiddleware)
	handler.RegisterReviewHandler(r, reviewService, preliminaryInvestigationService, authMiddleware)
	handler.RegisterPreliminaryInvestigationHandler(r, reviewService, preliminaryInvestigationService, authMiddleware)

	slog.Info("routes registered")

	r.Run(os.Getenv("PORT"))
}
