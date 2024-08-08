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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)
	db, err := db.NewPostgresDB()

	defer db.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	handler := handler.NewHandlerTenant()
	r.POST("/tenant", handler.CreateTenant)

	log.Fatalln(r.Run(":8080"))
}
