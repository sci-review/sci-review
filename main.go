package main

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	slog.Info("Starting application")

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Environment variables loaded")

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

	db, err := connectDbWithRetries(5, 7*time.Second)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = execMigrations(db)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	appCache := cacheInit()
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	refreshTokenRepo := repo.NewRefreshTokenRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)
	organizationRepo := repo.NewOrganizationRepo(db)
	organizationService := service.NewOrganizationService(organizationRepo)
	reviewRepoSql := repo.NewReviewRepoSql(db)
	reviewRepoCache := cacheDecorator.NewReviewRepoCache(reviewRepoSql, appCache)
	reviewService := service.NewReviewService(reviewRepoCache, userRepo)
	investigationRepoSql := repo.NewInvestigationRepoSql(db)
	investigationRepoCache := cacheDecorator.NewInvestigationRepoCache(investigationRepoSql, appCache)
	investigationService := service.NewInvestigationService(investigationRepoCache)
	slog.Info("services initialized")

	createAdminUser(userService)
	slog.Info("admin user created")

	tokenMiddleware := middleware.TokenMiddleware()
	authMiddleware := handler.AuthMiddleware()
	adminMiddleware := handler.AdminMiddleware()
	reviewMiddleware := middleware.ReviewMiddleware(reviewService)
	investigationMiddleware := middleware.InvestigationMiddleware(investigationService)
	slog.Info("middleware initialized")

	r := gin.Default()
	templateConfig(r)
	staticFilesConfig(r)
	configCors(r)

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))

	handler.RegisterHomeHandler(r, authMiddleware)
	handler.RegisterAuthHandler(r, authService)
	handler.RegisterUserHandler(r, userService)
	handler.RegisterAdminHandler(r, userService, authMiddleware, adminMiddleware)
	handler.RegisterOrganizationHandler(r, organizationService, authMiddleware)
	handler.RegisterReviewHandler(r, reviewService, investigationService, tokenMiddleware, reviewMiddleware, investigationMiddleware)
	handler.RegisterInvestigationHandler(r, reviewService, investigationService, tokenMiddleware, reviewMiddleware, investigationMiddleware)

	slog.Info("routes registered")

	r.Run(os.Getenv("PORT"))
}

func configCors(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))
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

func staticFilesConfig(r *gin.Engine) {
	slog.Info("Configuring static files")
	r.Static("/assets", "./assets")
	r.StaticFile("/android-chrome-192x192.png", "./assets/android-chrome-192x192.png")
	r.StaticFile("/android-chrome-512x512.png", "./assets/android-chrome-512x512.png")
	r.StaticFile("/apple-touch-icon.png", "./assets/apple-touch-icon.png")
	r.StaticFile("/favicon-16x16.png", "./assets/favicon-16x16.png")
	r.StaticFile("/favicon-32x32.png", "./assets/favicon-32x32.png")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.StaticFile("/site.webmanifest", "./assets/site.webmanifest")
	slog.Info("Static files configured")
}

func templateConfig(r *gin.Engine) {
	slog.Info("Configuring templates")
	r.LoadHTMLGlob("templates/**/*")
	slog.Info("Templates configured")
}

func connectDbWithRetries(maxRetries int, retryInterval time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	dataSourceName := os.Getenv("DATABASE_URL")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sqlx.Connect("pgx", dataSourceName)
		if err == nil {
			slog.Info("Database connected established")
			return db, nil
		}

		slog.Error("Failed to connect to the database", "attempt", attempt, "maxRetries", maxRetries, "error", err.Error())
		if attempt < maxRetries {
			time.Sleep(retryInterval)
		}
	}
	slog.Error("Failed to connect to the database after multiple attempts")

	return nil, errors.New("failed to connect to the database after multiple attempts")
}

func execMigrations(db *sqlx.DB) error {
	slog.Info("Migrating database")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Database migrated")
	return nil
}

func createAdminUser(userService *service.UserService) {
	slog.Info("Creating admin user")
	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	err := userService.CreateAdminUser(adminName, adminEmail, adminPassword)
	if err != nil {
		slog.Error("Error creating admin user", "error", err.Error())
		return
	}
}
