package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sci-review/auth"
	"sci-review/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dataSourceName := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo)

	r := gin.Default()
	auth.Register(r, authService)
	r.POST("/users", userHandler.Create)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(os.Getenv("PORT"))
}
