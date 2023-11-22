package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	PORT   string
	USER   string
	PASS   string
	DBNAME string
)

func CarrConfig() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		os.Exit(1)
	}

	USER = os.Getenv("USER")
	PASS = os.Getenv("PASS")
	PORT = os.Getenv("PORT")
	DBNAME = os.Getenv("DBNAME")
}
