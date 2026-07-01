package config

import (
	"os"

	"go.uber.org/zap"
)

var JwtSecret []byte

func init() {
	jwtEnv := os.Getenv("JWT_SECRET")
	if jwtEnv != "" {
		JwtSecret = []byte(jwtEnv)
	}

	if string(JwtSecret) == "minha_chave_super_secreta" {
		zap.L().Fatal("erro ao inicializar JwtSecret", zap.String("reason", "chave padrao utilizada em producao"))
	}

}
