package config

import (
	"log"
	"os"
)

var JwtSecret []byte

func init() {
	jwtEnv := os.Getenv("JWT_SECRET")
	if jwtEnv != "" {
		JwtSecret = []byte(jwtEnv)
	}

	if JwtSecret == nil {
		log.Fatal("erro ao inicializar JwtSecret")
	}

}
