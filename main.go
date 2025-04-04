package main

import (
	"api-enviar-email-moleniuk/config"
	"api-enviar-email-moleniuk/database"
	"api-enviar-email-moleniuk/handlers"
	"api-enviar-email-moleniuk/middleware"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

func initRedis() {
	_ = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func rateLimitMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Minute), 100)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func checkAndCreateUser(username, password, apiKey string) {
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Erro ao verificar usuário: %v", err)
		return
	}

	if exists == 0 {
		err = database.CreateUser(username, password, apiKey)
		if err != nil {
			log.Printf("Erro ao criar usuário %s: %v", username, err)
			return
		}
		fmt.Printf("Usuário %s criado com sucesso!\n", username)
	} else {
		fmt.Printf("Usuário %s existente.\n", username)
	}
}

func main() {
	config.LoadEnv()
	database.InitDB()
	initRedis()

	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASSWORD")
	apiKey := os.Getenv("API_KEY")

	if username == "" || password == "" || apiKey == "" {
		log.Println("Aviso: Variáveis de ambiente API_USERNAME, API_PASSWORD e API_KEY devem estar definidas.")
	} else {
		checkAndCreateUser(username, password, apiKey)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery(), rateLimitMiddleware())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-KEY")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Erro ao definir proxies confiáveis: %v", err)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/send-email", handlers.SendEmail)
	}

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
