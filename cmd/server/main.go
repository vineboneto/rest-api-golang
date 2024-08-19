package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/vineboneto/rest-api-golang/docs"

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

// @title						Crud Golang
// @version					1.0

// @description				Controle de Estoque com Go.
// @termsOfService				http://swagger.io/terms/

// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io

// @securityDefinitions.apiKey	JWT
// @in							header
// @name						token

// @license.name				MIT

// @host						localhost:8080
// @BasePath  /api/v1
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

	v1Router := r.Group("/api/v1")
	{
		v1Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		router := v1Router.Group("/tenant")
		{
			h := handler.NewHandlerTenant(db)
			router.POST("/", wrapper(h.CreateTenant))
			router.GET("/", wrapper(h.LoadAll))
			router.GET("/:id", wrapper(h.LoadById))
			router.PATCH("/:id", wrapper(h.Update))
		}

	}

	log.Fatalln(r.Run(":8080"))
}
