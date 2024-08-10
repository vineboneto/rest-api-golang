package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	handler "github.com/vineboneto/rest-api-golang/internal/handlers"
	"github.com/vineboneto/rest-api-golang/internal/infra/db"
)

type GinHandler func(c *gin.Context) error

func wrapper(handler GinHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)
	db, err := db.NewPostgresDB()

	defer func() {
		db.Close()
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router := r.Group("/tenant")
	{
		handler := handler.NewHandlerTenant(db)
		router.POST("/", wrapper(handler.CreateTenant))
		router.GET("/", wrapper(handler.LoadAll))
		router.GET("/:id", wrapper(handler.LoadById))
		router.PATCH("/:id", wrapper(handler.Update))
	}

	log.Fatalln(r.Run(":8080"))
}
