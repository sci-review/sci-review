package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
	"os"
	cacheDecorator "sci-review/cache"
	"sci-review/handler"
	"sci-review/middleware"
	"sci-review/repo"
	"sci-review/service"
	"strconv"
	"time"
)

func main() {
	slog.Info("starting application")

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("environment variables loaded")

	appEnv := os.Getenv("APP_ENV")
	logLevel := slog.LevelDebug
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		logLevel = slog.LevelInfo
	}
	slog.Info("set application environment", "appEnv", appEnv)

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	slog.Info("logger initialized")

	dataSourceName := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("database connected established")

	appCache := cacheInit()
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo)
	organizationRepo := repo.NewOrganizationRepo(db)
	organizationService := service.NewOrganizationService(organizationRepo)
	reviewRepo := repo.NewReviewRepo(db, appCache)
	reviewService := service.NewReviewService(reviewRepo)
	investigationRepoSql := repo.NewInvestigationRepoSql(db)
	investigationRepoCache := cacheDecorator.NewInvestigationRepoCache(investigationRepoSql, appCache)
	investigationService := service.NewInvestigationService(investigationRepoCache)
	slog.Info("services initialized")

	authMiddleware := handler.AuthMiddleware()
	adminMiddleware := handler.AdminMiddleware()
	reviewMiddleware := middleware.ReviewMiddleware(reviewService)
	investigationMiddleware := middleware.InvestigationMiddleware(investigationService)
	slog.Info("middleware initialized")

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/assets", "./assets")

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))

	handler.RegisterHomeHandler(r, authMiddleware)
	handler.RegisterAuthHandler(r, authService)
	handler.RegisterUserHandler(r, userService)
	handler.RegisterAdminHandler(r, userService, authMiddleware, adminMiddleware)
	handler.RegisterOrganizationHandler(r, organizationService, authMiddleware)
	handler.RegisterReviewHandler(r, reviewService, investigationService, authMiddleware, reviewMiddleware, investigationMiddleware)
	handler.RegisterInvestigationHandler(r, reviewService, investigationService, authMiddleware, reviewMiddleware, investigationMiddleware)

	slog.Info("routes registered")

	r.Run(os.Getenv("PORT"))
}

func cacheInit() *cache.Cache {
	slog.Info("Initializing cache")
	defaultExpirationMinutesStr := os.Getenv("CACHE_DEFAULT_EXPIRATION_MINUTES")
	cleanupIntervalMinutesStr := os.Getenv("CACHE_CLEANUP_INTERVAL_MINUTES")

	defaultExpirationMinutes, err := strconv.Atoi(defaultExpirationMinutesStr)
	if err != nil {
		slog.Error("Error converting cache default expiration to int", "error", err.Error())
		return nil
	}

	cleanupIntervalMinutes, err := strconv.Atoi(cleanupIntervalMinutesStr)
	if err != nil {
		slog.Error("Error converting cache cleanup interval to int", "error", err.Error())
		return nil
	}
	defaultExpiration := time.Duration(defaultExpirationMinutes) * time.Minute
	cleanupInterval := time.Duration(cleanupIntervalMinutes) * time.Minute

	slog.Info("Cache initialized", "defaultExpiration in minutes", defaultExpirationMinutes, "cleanupInterval in minutes", cleanupIntervalMinutes)
	return cache.New(defaultExpiration, cleanupInterval)
}
