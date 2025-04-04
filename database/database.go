package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Erro: Todas as variáveis de ambiente DB_HOST, DB_PORT, DB_USER, DB_PASSWORD e DB_NAME devem estar definidas.")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao pingar o PostgreSQL: %v", err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			api_key TEXT UNIQUE NOT NULL
		);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	log.Println("Conexão com PostgreSQL estabelecida e tabela criada com sucesso!")
}

func CreateUser(username, password, apiKey string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}
	_, err = DB.Exec("INSERT INTO users (username, password, api_key) VALUES ($1, $2, $3)", username, hashedPassword, apiKey)
	return err
}

func ValidateUser(username, password string) bool {
	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) == nil
}

func ValidateAPIKey(apiKey string) bool {
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE api_key = $1", apiKey).Scan(&exists)
	return err == nil && exists > 0
}
