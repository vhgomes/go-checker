package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserID extrai o user_id do contexto Gin de forma segura.
// O MiddlewareJWT já garante que o valor é uint — esta função apenas encapsula
// o acesso e retorna um erro HTTP padronizado se a chave não existir.
func GetUserID(c *gin.Context) (uint, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autorizado"})
		return 0, false
	}

	id, ok := v.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id com tipo inválido no contexto"})
		return 0, false
	}

	return id, true
}