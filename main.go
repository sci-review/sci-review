package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"
	"sci-review/user"
)

func main() {
	dataSourceName := "postgresql://postgres:postgres@localhost:5432/sci_review"
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	r := gin.Default()
	r.POST("/users", userHandler.Create)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
